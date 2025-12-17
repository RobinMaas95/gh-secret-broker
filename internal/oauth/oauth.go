package oauth

import (
	"encoding/gob"
	"encoding/json"
	"log/slog"
	"net/http"
	"sort"

	"github.com/RobinMaas95/gh-secret-broker/internal/config"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
)

func init() {
	gob.Register(goth.User{})
}

type ProviderIndex struct {
	Providers    []string
	ProvidersMap map[string]string
}

type Service struct {
	logger *slog.Logger
	config *config.Config
	store  sessions.Store
}

func NewService(logger *slog.Logger, cfg *config.Config) *Service {
	// Initialize the store with secret from config
	store := sessions.NewFilesystemStore("", []byte(cfg.SessionSecret))
	store.MaxLength(0) // No limit on length
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = cfg.IsProduction()
	store.Options.MaxAge = 86400 * 1 // 1 day

	gothic.Store = store

	gothic.GetProviderName = func(req *http.Request) (string, error) {
		return req.PathValue("provider"), nil
	}

	if cfg.GithubEnterpriseURL != "" {
		github.AuthURL = cfg.GithubEnterpriseURL + "/login/oauth/authorize"
		github.TokenURL = cfg.GithubEnterpriseURL + "/login/oauth/access_token"
		github.ProfileURL = cfg.GithubEnterpriseURL + "/api/v3/user"
		github.EmailURL = cfg.GithubEnterpriseURL + "/api/v3/user/emails"
	}

	goth.UseProviders(
		github.New(cfg.GithubClientID, cfg.GithubClientSecret, cfg.BaseURL+"/auth/github/callback", "user:email"),
	)

	return &Service{
		logger: logger,
		config: cfg,
		store:  store,
	}
}

func (s *Service) GetProviderIndex() *ProviderIndex {
	m := map[string]string{
		"github": "Github",
	}
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	return &ProviderIndex{Providers: keys, ProvidersMap: m}
}

func (s *Service) ProviderLogin(res http.ResponseWriter, req *http.Request) {
	// try to get the user without re-authenticating
	session, err := s.store.Get(req, "session")
	if err != nil {
		s.logger.Warn("Failed to get session, creating new one", slog.String("error", err.Error()))
		// Continue with new/invalid session - don't block login
	}

	if session != nil {
		if _, ok := session.Values["user"]; ok {
			// User is already logged in, redirect to user page
			http.Redirect(res, req, "/", http.StatusTemporaryRedirect)
			return
		}
	}

	// Fallback to gothic check (optional, but good for consistency if mixed)
	if _, err := gothic.CompleteUserAuth(res, req); err == nil {
		http.Redirect(res, req, "/", http.StatusTemporaryRedirect)
	} else {
		gothic.BeginAuthHandler(res, req)
	}
}

func (s *Service) ProviderLogout(res http.ResponseWriter, req *http.Request) {
	if err := gothic.Logout(res, req); err != nil {
		s.logger.Warn("Gothic logout error", slog.String("error", err.Error()))
	}

	session, err := s.store.Get(req, "session")
	if err != nil {
		s.logger.Warn("Failed to get session during logout", slog.String("error", err.Error()))
	}

	if session != nil {
		delete(session.Values, "user")
		if err := session.Save(req, res); err != nil {
			s.logger.Error("Failed to save session during logout", slog.String("error", err.Error()))
		}
	}

	res.Header().Set("Location", "/login")
	res.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	res.WriteHeader(http.StatusTemporaryRedirect)
}

func (s *Service) HandleCallback(res http.ResponseWriter, req *http.Request) {
	user, err := gothic.CompleteUserAuth(res, req)
	if err != nil {
		s.logger.Error("Failed to complete user auth", slog.String("error", err.Error()))
		http.Error(res, "Authentication failed", http.StatusInternalServerError)
		return
	}
	s.logger.Info("User logged in", slog.String("user_id", user.UserID), slog.String("email", user.Email))

	// Store user in session
	session, err := s.store.Get(req, "session")
	if err != nil {
		s.logger.Error("Failed to get session after auth", slog.String("error", err.Error()))
		http.Error(res, "Session error", http.StatusInternalServerError)
		return
	}

	session.Values["user"] = user
	if err = session.Save(req, res); err != nil {
		s.logger.Error("Failed to save session", slog.String("error", err.Error()))
		http.Error(res, "Failed to save session", http.StatusInternalServerError)
		return
	}

	// Redirect to the user page after successful login
	http.Redirect(res, req, "/", http.StatusTemporaryRedirect)
}

/*
HandleProvidersAPI returns the list of available providers.
It is used by the frontend to display the list of providers that might
change over time.
*/
func (s *Service) HandleProvidersAPI(res http.ResponseWriter, req *http.Request) {
	index := s.GetProviderIndex()
	res.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(res).Encode(index); err != nil {
		s.logger.Error("Failed to encode providers response", slog.String("error", err.Error()))
	}
}

/*
HandleUserAPI returns the user information from the session.
It is used by the frontend to display the user information.
*/
func (s *Service) HandleUserAPI(res http.ResponseWriter, req *http.Request) {
	session, err := s.store.Get(req, "session")
	if err != nil {
		s.logger.Warn("HandleUserAPI: Failed to get session", slog.String("error", err.Error()))
		http.Error(res, "Unauthorized", http.StatusUnauthorized)
		return
	}

	val, ok := session.Values["user"]
	if !ok {
		s.logger.Debug("HandleUserAPI: No user in session")
		http.Error(res, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, ok := val.(goth.User)
	if !ok {
		s.logger.Error("HandleUserAPI: User in session is not goth.User")
		http.Error(res, "Unauthorized", http.StatusUnauthorized)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(res).Encode(user); err != nil {
		s.logger.Error("Failed to encode user response", slog.String("error", err.Error()))
	}
}
