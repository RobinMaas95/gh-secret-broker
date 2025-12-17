package main

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"io/fs"

	"github.com/RobinMaas95/gh-secret-broker/ui"
	"github.com/justinas/nosurf"
)

func ping(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("OK")) //nolint:errcheck // health check endpoint
}

func (app *application) handleCsrfToken(w http.ResponseWriter, r *http.Request) {
	token := nosurf.Token(r)
	w.Header().Set("Content-Type", "application/json")
	// We can reuse the same struct or map, keep it simple for now
	if err := json.NewEncoder(w).Encode(map[string]string{"token": token}); err != nil {
		app.logger.Error("Failed to encode CSRF token", slog.String("error", err.Error()))
	}
}

func (app *application) handleSPA(w http.ResponseWriter, r *http.Request) {
	distFS, err := fs.Sub(ui.Files, "build")
	if err != nil {
		app.logger.Error("Could not get static files", slog.String("error", err.Error()))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Try to serve the file directly -> Not part of the SPA
	path := r.URL.Path
	if path == "/" {
		path = "index.html"
	} else if len(path) > 0 && path[0] == '/' {
		path = path[1:]
	}

	f, err := distFS.Open(path)
	if err == nil {
		defer func() { _ = f.Close() }()
		stats, _ := f.Stat()
		if !stats.IsDir() {
			http.FileServer(http.FS(distFS)).ServeHTTP(w, r)
			return
		}
	}

	// Fallback to index.html for non-file requests (SPA routing)
	content, err := fs.ReadFile(distFS, "index.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	_, _ = w.Write(content)
}
