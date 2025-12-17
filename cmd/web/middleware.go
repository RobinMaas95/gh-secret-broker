package main

import (
	"fmt"
	"net/http"

	"github.com/justinas/nosurf"
)

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			// We use the built-in recover() to check if an error occurred. If so,
			// recover returns a panic-value
			pv := recover()
			if pv != nil {
				w.Header().Set("Connection", "close")
				app.serverError(w, r, fmt.Errorf("%v", pv))

			}
		}()
		next.ServeHTTP(w, r)
	})
}

// preventCSRF is a factory function that creates CSRF protection middleware.
// It needs to be called with the application config to set cookie security properly.
// The factory pattern is used because a middleware has a specific signature (func(http.Handler) http.Handler)
// meaning we cannot pass in parameters. Using the factory, we can pass in the parameter and access it in the middleware.
func preventCSRFFactory(isProduction bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		csrfHandler := nosurf.New(next)
		csrfHandler.SetBaseCookie(http.Cookie{
			HttpOnly: true,
			Path:     "/",
			Secure:   isProduction, // Only send over HTTPS in production
			SameSite: http.SameSiteLaxMode,
		})
		return csrfHandler
	}
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			ip     = r.RemoteAddr
			proto  = r.Proto
			method = r.Method
			uri    = r.URL.RequestURI()
		)

		app.logger.Info("received request", "ip", ip, "proto", proto, "method", method, "uri", uri)

		next.ServeHTTP(w, r)
	})
}

func (app *application) commonHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' 'sha256-OWM7eqRFZQb+Ir51EYZz6GAUUG4Oy2fNwCWX0peasbE=' 'sha256-tcbDxjMo+xKqM21aCGYbs/QAJqB7yUXC06oPWDapBgc=' 'sha256-zWpgYAIYQbPPXWm2cNN92poH5pezyiyARDiGUjuqbFU=' 'sha256-NDlUvbI0C5AhCY+uu2OxERc8b/zOZ5m/C3vpWbghG1M=' 'sha256-bo9/JAqIUBMiSHL1O4oiO3U5UHaFxqbagFBryI+8mwU=' 'unsafe-hashes' fonts.googleapis.com; script-src 'self' 'unsafe-eval' 'unsafe-inline'; font-src fonts.gstatic.com; img-src 'self' avatars.githubusercontent.com")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")
		w.Header().Set("Server", "Go")

		next.ServeHTTP(w, r)
	})
}
