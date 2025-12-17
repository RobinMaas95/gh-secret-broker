package main

import (
	"net/http"

	"github.com/RobinMaas95/gh-secret-broker/internal/oauth"
	"github.com/justinas/alice"
)

func (app *application) routes(oauthService *oauth.Service) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", app.handleSPA)
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
	mux.Handle("PUT /api/repo/{owner}/{repo}/secrets/{name}", dynamic.ThenFunc(app.handleCreateSecret))

	standard := alice.New(app.recoverPanic, app.logRequest, app.commonHeaders)
	return standard.Then(mux)
}
