package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
)

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
		trace  = string(debug.Stack())
	)

	text := http.StatusText(http.StatusInternalServerError)
	if app.debugMode {
		text = fmt.Sprintf("%s\n%s", err, trace)
	}
	app.logger.Error(err.Error(), slog.String("method", method), slog.String("uri", uri), slog.String("trace", trace))
	http.Error(w, text, http.StatusInternalServerError)
}
