package main

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/justinas/nosurf"
)

func TestRecoverPanic(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	app := &application{logger: logger}

	panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	})

	wrapped := app.recoverPanic(panicHandler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	wrapped.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected 500, got %d", w.Code)
	}
	if w.Header().Get("Connection") != "close" {
		t.Error("Expected Connection: close header")
	}
}

func TestLogRequest(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo}))
	app := &application{logger: logger}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	wrapped := app.logRequest(handler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	wrapped.ServeHTTP(w, req)

	logOutput := buf.String()
	if !strings.Contains(logOutput, "received request") {
		t.Error("Expected 'received request' in logs")
	}
	if !strings.Contains(logOutput, "GET") {
		t.Error("Expected 'GET' in logs")
	}
	if !strings.Contains(logOutput, "/") {
		t.Error("Expected '/' in logs")
	}
}

func TestCommonHeaders(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	app := &application{logger: logger}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	wrapped := app.commonHeaders(handler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	wrapped.ServeHTTP(w, req)

	headers := w.Header()
	if headers.Get("Content-Security-Policy") == "" {
		t.Error("Expected Content-Security-Policy header")
	}
	if headers.Get("X-Frame-Options") != "deny" {
		t.Error("Expected X-Frame-Options: deny")
	}
	if headers.Get("X-Content-Type-Options") != "nosniff" {
		t.Error("Expected X-Content-Type-Options: nosniff")
	}
	if headers.Get("Referrer-Policy") != "origin-when-cross-origin" {
		t.Error("Expected Referrer-Policy: origin-when-cross-origin")
	}
	if headers.Get("Server") != "Go" {
		t.Error("Expected Server: Go")
	}
}

func TestPreventCSRF(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nosurf.Token(r) // Force token generation to ensure cookie is set
		w.WriteHeader(http.StatusOK)
	})

	// Test both development and production modes
	for _, isProduction := range []bool{false, true} {
		wrapped := preventCSRFFactory(isProduction)(handler)
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()
		wrapped.ServeHTTP(w, req)

		// CSRF middleware sets cookie on first request
		cookies := w.Result().Cookies()
		csrfFound := false
		for _, cookie := range cookies {
			if cookie.Name == "csrf_token" {
				csrfFound = true
				if cookie.Secure != isProduction {
					t.Errorf("Expected Secure=%v for isProduction=%v, got %v", isProduction, isProduction, cookie.Secure)
				}
				break
			}
		}
		if !csrfFound {
			t.Error("Expected CSRF cookie")
		}
	}
}
