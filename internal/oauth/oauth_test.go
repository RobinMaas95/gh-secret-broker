package oauth

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RobinMaas95/gh-secret-broker/internal/config"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
)

func TestGetProviderIndex(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	cfg := &config.Config{SessionSecret: "testsecret"}
	svc := NewService(logger, cfg)

	index := svc.GetProviderIndex()
	if len(index.Providers) != 1 || index.Providers[0] != "github" {
		t.Errorf("Expected github provider, got %v", index.Providers)
	}
	if index.ProvidersMap["github"] != "Github" {
		t.Errorf("Expected Github name, got %q", index.ProvidersMap["github"])
	}
}

func TestHandleProvidersAPI(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	cfg := &config.Config{SessionSecret: "testsecret"}
	svc := NewService(logger, cfg)

	handler := http.HandlerFunc(svc.HandleProvidersAPI)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Expected JSON content-type, got %q", w.Header().Get("Content-Type"))
	}

	var index ProviderIndex
	if err := json.NewDecoder(w.Body).Decode(&index); err != nil {
		t.Fatalf("Failed to decode JSON: %v", err)
	}

	if len(index.Providers) != 1 || index.Providers[0] != "github" {
		t.Errorf("Expected github provider, got %v", index.Providers)
	}
}

func TestHandleUserAPI_NoUser(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	cfg := &config.Config{SessionSecret: "testsecret"}
	svc := NewService(logger, cfg)

	handler := http.HandlerFunc(svc.HandleUserAPI)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestProviderLogin_NoUser(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	cfg := &config.Config{SessionSecret: "testsecret"}
	svc := NewService(logger, cfg)

	handler := http.HandlerFunc(svc.ProviderLogin)
	req, err := http.NewRequest(http.MethodGet, "/auth/github?provider=github", nil)
	req.SetPathValue("provider", "github")
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	// Should call gothic.BeginAuthHandler, which redirects to GitHub
	// Goth expects the provider name in the URL path
	if w.Code != http.StatusTemporaryRedirect {
		t.Errorf("Expected redirect to auth (302), got status %d (body: %s)", w.Code, w.Body.String())
	}
}

func TestProviderLogin_UserLoggedIn(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	cfg := &config.Config{SessionSecret: "testsecret"}

	// Use CookieStore for testing to easily inject session data
	store := sessions.NewCookieStore([]byte(cfg.SessionSecret))

	svc := &Service{
		logger: logger,
		config: cfg,
		store:  store,
	}

	req, _ := http.NewRequest(http.MethodGet, "/auth/github?provider=github", nil)
	w := httptest.NewRecorder()

	// Pre-populate session
	session, _ := store.Get(req, "session")
	session.Values["user"] = goth.User{UserID: "123", Email: "test@example.com"}
	_ = session.Save(req, w)

	// Create a new request with the cookie set
	req, _ = http.NewRequest(http.MethodGet, "/auth/github?provider=github", nil)
	req.Header.Set("Cookie", w.Header().Get("Set-Cookie"))
	w = httptest.NewRecorder()

	handler := http.HandlerFunc(svc.ProviderLogin)
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusTemporaryRedirect {
		t.Errorf("Expected redirect (307), got %d", w.Code)
	}
	if loc := w.Header().Get("Location"); loc != "/#/userpage" {
		t.Errorf("Expected redirect to /#/userpage, got %s", loc)
	}
}

func TestProviderLogout(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	cfg := &config.Config{SessionSecret: "testsecret"}
	store := sessions.NewCookieStore([]byte(cfg.SessionSecret))

	svc := &Service{
		logger: logger,
		config: cfg,
		store:  store,
	}

	req, _ := http.NewRequest(http.MethodGet, "/logout/github", nil)
	w := httptest.NewRecorder()

	// Pre-populate session
	session, _ := store.Get(req, "session")
	session.Values["user"] = goth.User{UserID: "123"}
	_ = session.Save(req, w)

	// Request with cookie
	req.Header.Set("Cookie", w.Header().Get("Set-Cookie"))
	w = httptest.NewRecorder()

	handler := http.HandlerFunc(svc.ProviderLogout)
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusTemporaryRedirect {
		t.Errorf("Expected redirect (307), got %d", w.Code)
	}
	if loc := w.Header().Get("Location"); loc != "/" {
		t.Errorf("Expected redirect to /, got %s", loc)
	}

	// Verify user is removed from session
	res := w.Result()
	cookies := res.Cookies()
	var sessionCookie *http.Cookie
	for _, c := range cookies {
		if c.Name == "session" {
			sessionCookie = c
			break
		}
	}

	if sessionCookie != nil {
		// Create a request with this new cookie to check its contents using the store
		checkReq, _ := http.NewRequest("GET", "/", nil)
		checkReq.AddCookie(sessionCookie)
		checkSession, _ := store.Get(checkReq, "session")
		if _, ok := checkSession.Values["user"]; ok {
			t.Error("User should be removed from session")
		}
	}
}

func TestHandleUserAPI_Authenticated(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	cfg := &config.Config{SessionSecret: "testsecret"}
	store := sessions.NewCookieStore([]byte(cfg.SessionSecret))

	svc := &Service{
		logger: logger,
		config: cfg,
		store:  store,
	}

	req, _ := http.NewRequest(http.MethodGet, "/api/me", nil)
	w := httptest.NewRecorder()

	// Pre-populate session
	session, _ := store.Get(req, "session")
	expectedUser := goth.User{UserID: "123", Email: "test@example.com", Name: "Test User"}
	session.Values["user"] = expectedUser
	_ = session.Save(req, w)

	// Request with cookie
	req.Header.Set("Cookie", w.Header().Get("Set-Cookie"))
	w = httptest.NewRecorder()

	handler := http.HandlerFunc(svc.HandleUserAPI)
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var gotUser goth.User
	if err := json.NewDecoder(w.Body).Decode(&gotUser); err != nil {
		t.Fatalf("Failed to decode JSON: %v", err)
	}

	if gotUser.UserID != expectedUser.UserID {
		t.Errorf("Expected UserID %s, got %s", expectedUser.UserID, gotUser.UserID)
	}
}
