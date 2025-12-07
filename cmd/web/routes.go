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

	dynamic := alice.New(preventCSRF)
	mux.HandleFunc("GET /auth/{provider}/callback", oauthService.HandleCallback)
	mux.Handle("GET /logout/{provider}", dynamic.ThenFunc(oauthService.ProviderLogout))
	mux.Handle("GET /auth/{provider}", dynamic.ThenFunc(oauthService.ProviderLogin))

	mux.HandleFunc("GET /api/providers", oauthService.HandleProvidersAPI)
	mux.HandleFunc("GET /api/me", oauthService.HandleUserAPI)
	mux.Handle("GET /api/repositories", dynamic.ThenFunc(app.handleListRepositories))

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
