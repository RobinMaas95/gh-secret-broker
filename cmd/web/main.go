package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"slices"
	"time"

	"github.com/RobinMaas95/gh-secret-broker/internal/config"
	"github.com/RobinMaas95/gh-secret-broker/ui"
	"github.com/joho/godotenv"
)

var logHandler slog.Handler

func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got / request\n")
	_, err := io.WriteString(w, "This is my website!\n")
	if err != nil {
		fmt.Println("Could not write to stream")
		os.Exit(1)
	}
}

func setupLogger(logFormat string) slog.Handler {
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	switch logFormat {
	case "json":
		logHandler = slog.NewJSONHandler(os.Stdout, opts)
	case "text":
		opts.AddSource = true
		logHandler = slog.NewTextHandler(os.Stdout, opts)
	default:
		// Should never be reached because we validate the user input
		opts.AddSource = true
		logHandler = slog.NewTextHandler(os.Stdout, opts)
	}
	return logHandler
}

func readEnvs() {
}

func main() {
	logFormat := "text"
	_ = godotenv.Load() // Load .env file if it exists
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

	logger := slog.New(setupLogger(logFormat))
	slog.SetDefault(logger)

	mux := http.NewServeMux()
	distFS, err := fs.Sub(ui.Files, "dist")
	if err != nil {
		logger.Error("Could not get dist files", slog.String("error", err.Error()))
		os.Exit(1)
	}
	mux.Handle("GET /", http.FileServer(http.FS(distFS)))

	srv := &http.Server{
		Addr:         cfg.Addr,
		Handler:      mux,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	logger.Info("Starting server", slog.String("addr", srv.Addr))
	err = srv.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}
