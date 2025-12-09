package main

import (
	"encoding/json"
	"log/slog"
	"net/http"

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
