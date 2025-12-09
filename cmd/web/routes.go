package main

import (
	"io/fs"
	"log/slog"
	"net/http"

	"github.com/RobinMaas95/gh-secret-broker/internal/oauth"
	"github.com/RobinMaas95/gh-secret-broker/ui"
	"github.com/justinas/alice"
)

func (app *application) routes(oauthService *oauth.Service) http.Handler {
	mux := http.NewServeMux()
	distFS, err := fs.Sub(ui.Files, "dist")
	if err != nil {
		app.logger.Error("Could not get dist files", slog.String("error", err.Error()))
		panic("failed to get dist files: " + err.Error())
	}
	mux.Handle("GET /", http.FileServer(http.FS(distFS)))
	mux.HandleFunc("GET /ping", ping)

	dynamic := alice.New(preventCSRFFactory(app.config.IsProduction()))

	// CSRF Protection Strategy:
	// - OAuth flows (login/logout) use CSRF protection via the 'dynamic' middleware
	// - All current API endpoints are GET requests (read-only) and don't need CSRF
	// - SameSite=Lax cookies provide additional protection against cross-site requests
	//
	// IMPORTANT: When adding POST/PUT/DELETE endpoints in the future:
	// Apply the 'dynamic' middleware but ensure your SPA sends CSRF tokens
	//
	// Example for future state-changing endpoints:
	//   mux.Handle("POST /api/secrets", dynamic.ThenFunc(app.handleCreateSecret))

	mux.HandleFunc("GET /auth/{provider}/callback", oauthService.HandleCallback)
	mux.Handle("GET /logout/{provider}", dynamic.ThenFunc(oauthService.ProviderLogout))
	mux.Handle("GET /auth/{provider}", dynamic.ThenFunc(oauthService.ProviderLogin))

	// API Routes - All currently read-only (GET)
	// No CSRF needed as they don't change state and SameSite cookies provide protection
	mux.HandleFunc("GET /api/providers", oauthService.HandleProvidersAPI)
	mux.HandleFunc("GET /api/user", oauthService.HandleUserAPI)
	mux.HandleFunc("GET /api/user/repos", app.handleListRepositories)
	mux.HandleFunc("GET /api/repo/{owner}/{repo}/secrets", app.handleListSecrets)
	// Dynamic because we need our preventCSRFFactory to be applied so that
	// our token endpoint returns a valid token
	mux.Handle("GET /api/csrf-token", dynamic.ThenFunc(app.handleCsrfToken))
	mux.Handle("DELETE /api/repo/{owner}/{repo}/secrets/{name}", dynamic.ThenFunc(app.handleDeleteSecret))

	// Serve index.html for /userpage to support SPA history mode (if used)
	mux.HandleFunc("GET /userpage", func(w http.ResponseWriter, r *http.Request) {
		content, err := fs.ReadFile(distFS, "index.html")
		if err != nil {
			app.logger.Error("Failed to read index.html", slog.String("error", err.Error()))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		if _, err := w.Write(content); err != nil {
			app.logger.Error("Failed to write response", slog.String("error", err.Error()))
		}
	})

	standard := alice.New(app.recoverPanic, app.logRequest, app.commonHeaders)
	return standard.Then(mux)
}
