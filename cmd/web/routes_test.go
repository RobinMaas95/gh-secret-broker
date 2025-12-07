package main

import (
	"log/slog"
	"testing"

	"github.com/RobinMaas95/gh-secret-broker/internal/config"
	"github.com/RobinMaas95/gh-secret-broker/internal/oauth"
)

func TestRoutes(t *testing.T) {
	// Setup minimal dependencies for the application
	cfg := &config.Config{
		SessionSecret: "test-secret",
	}
	app := &application{
		config: cfg,
		logger: slog.Default(),
	}

	// Mock OAuth service (dependencies can be nil/minimal as we just test route registration)
	oauthService := &oauth.Service{}

	// This function call will panic if there are conflicting routes
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("The application panicked during route registration: %v", r)
		}
	}()

	mux := app.routes(oauthService)

	if mux == nil {
		t.Fatal("Expected a router to be returned")
	}
}
