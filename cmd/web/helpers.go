package main

import (
	"log/slog"
	"net/http"

	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
)

// requireUser checks if the user is authenticated in the session.
// If valid, it returns the user and true.
// If invalid, it writes an Unauthorized error response and returns false.
func (app *application) requireUser(w http.ResponseWriter, r *http.Request) (goth.User, bool) {
	session, err := gothic.Store.Get(r, "session")
	if err != nil {
		app.logger.Error("Failed to get session", slog.String("error", err.Error()))
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return goth.User{}, false
	}

	val, ok := session.Values["user"]
	if !ok {
		app.logger.Debug("No user in session")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return goth.User{}, false
	}

	user, ok := val.(goth.User)
	if !ok {
		app.logger.Error("User in session is not goth.User")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return goth.User{}, false
	}

	return user, true
}

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Error("server error", slog.String("error", err.Error()), slog.String("method", r.Method), slog.String("uri", r.URL.RequestURI()))
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
