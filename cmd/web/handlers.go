package main

import "net/http"

func ping(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("OK")) //nolint:errcheck // health check endpoint
}
