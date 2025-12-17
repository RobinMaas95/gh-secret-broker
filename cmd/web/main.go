package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"slices"
	"syscall"
	"time"

	"github.com/RobinMaas95/gh-secret-broker/internal/config"
	"github.com/RobinMaas95/gh-secret-broker/internal/oauth"
	"github.com/RobinMaas95/gh-secret-broker/internal/repository"
	"github.com/google/go-github/v80/github"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

type application struct {
	logger       *slog.Logger
	debugMode    bool
	repositories repository.RepositoryService
	config       *config.Config
	patClient    *github.Client
}

func setupLogger(logFormat string) slog.Handler {
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	switch logFormat {
	case "json":
		return slog.NewJSONHandler(os.Stdout, opts)
	case "text":
		opts.AddSource = true
		return slog.NewTextHandler(os.Stdout, opts)
	default:
		// Should never be reached because we validate the user input
		opts.AddSource = true
		return slog.NewTextHandler(os.Stdout, opts)
	}
}

func main() {
	_ = godotenv.Load() // Load .env file if it exists

	logFormat := "text"
	cfg, err := config.Load()
	if err != nil {
		fmt.Println(err)
		fmt.Println("Could not load config")
		os.Exit(1)
	}
	flag.StringVar(&cfg.Addr, "addr", ":4000", "HTTP network address")
	flag.Func("log-format", "Output format of the logs", func(flagValue string) error {
		if slices.Contains([]string{"text", "json"}, flagValue) {
			logFormat = flagValue
			return nil
		}
		return errors.New(`must be one of "text" or "json"`)
	})
	flag.Parse() // Need to be called, so that the flag values are filled with passed values and not defaults
	fmt.Printf("The log format is: %s\n", logFormat)

	if cfg.BaseURL == "" {
		_, port, err := net.SplitHostPort(cfg.Addr)
		if err != nil {
			// If SplitHostPort fails (e.g. ":4000"), try to see if it's just a port or fallback
			if len(cfg.Addr) > 0 && cfg.Addr[0] == ':' {
				port = cfg.Addr[1:]
			} else {
				// Fallback: assume the whole thing is a port if it parses as int, otherwise default to 4000
				port = "4000" // Safe default
			}
		}
		cfg.BaseURL = "http://localhost:" + port
	}

	logger := slog.New(setupLogger(logFormat))
	slog.SetDefault(logger)

	// Initialize PAT Client
	// The pat client is used to access the GitHub API
	// on behalf of the application because the user's token is
	// not powerful enough to access secrets.
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: cfg.GithubPAT},
	)
	tc := oauth2.NewClient(context.Background(), ts)
	var patClient *github.Client

	if cfg.GithubEnterpriseURL != "" {
		baseURL := cfg.GithubEnterpriseURL + "/api/v3/"
		uploadURL := cfg.GithubEnterpriseURL + "/api/v3/upload/"
		var err error
		patClient, err = github.NewClient(tc).WithEnterpriseURLs(baseURL, uploadURL)
		if err != nil {
			logger.Error("Failed to create enterprise client", slog.String("error", err.Error()))
			os.Exit(1)
		}
	} else {
		patClient = github.NewClient(tc)
	}

	app := &application{
		logger:       logger,
		debugMode:    false,
		repositories: repository.NewService(),
		config:       cfg,
		patClient:    patClient,
	}

	oauthService := oauth.NewService(logger, cfg)

	srv := &http.Server{
		Addr:         cfg.Addr,
		Handler:      app.routes(oauthService),
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Graceful shutdown setup
	shutdownComplete := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)
		<-sigint

		logger.Info("Shutting down server...")

		// Give outstanding requests 30 seconds to complete
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			logger.Error("Server shutdown error", slog.String("error", err.Error()))
		}
		close(shutdownComplete)
	}()

	logger.Info("Starting server", slog.String("addr", srv.Addr))
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("Server error", slog.String("error", err.Error()))
		os.Exit(1)
	}

	<-shutdownComplete
	logger.Info("Server stopped gracefully")
}
