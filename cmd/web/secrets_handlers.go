package main

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/google/go-github/v80/github"
	"golang.org/x/oauth2"
)

func (app *application) handleListSecrets(w http.ResponseWriter, r *http.Request) {
	// Get Session & Verify User
	user, ok := app.requireUser(w, r)
	if !ok {
		return
	}

	// Extract path parameters
	owner := r.PathValue("owner")
	repo := r.PathValue("repo")

	if owner == "" || repo == "" {
		http.Error(w, "Missing owner or repo", http.StatusBadRequest)
		return
	}

	// We do not work on userGhClient directly, but instead pass it to the repository service.
	// In tests, we can mock the repository service and inject a mock client or just
	// return a fixed value.
	userGhClient, err := app.getGitHubClient(r.Context(), user.AccessToken)
	if err != nil {
		app.logger.Error("Failed to create GitHub client", slog.String("error", err.Error()))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	hasAccess, err := app.repositories.HasMaintainerAccess(r.Context(), userGhClient, owner, repo)

	if err != nil {
		app.logger.Error("Failed to check permissions", slog.String("error", err.Error()))
		http.Error(w, "Permission check failed", http.StatusInternalServerError)
		return
	}
	if !hasAccess {
		app.logger.Warn("User attempted to access secrets without permission", slog.String("user", user.Email), slog.String("repo", owner+"/"+repo))
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Use Shared PAT Client (Only after verification)
	githubClient := app.patClient

	// Call repository service
	secrets, err := app.repositories.ListSecrets(r.Context(), githubClient, owner, repo)
	if err != nil {
		app.logger.Error("Failed to list secrets", slog.String("error", err.Error()))
		http.Error(w, "Failed to list secrets", http.StatusInternalServerError)
		return
	}

	// Respond

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(secrets); err != nil {
		app.logger.Error("Failed to encode secrets", slog.String("error", err.Error()))
	}
}

func (app *application) handleDeleteSecret(w http.ResponseWriter, r *http.Request) {
	// Get Session & Verify User
	user, ok := app.requireUser(w, r)
	if !ok {
		return
	}

	// Extract path parameters
	owner := r.PathValue("owner")
	repo := r.PathValue("repo")
	name := r.PathValue("name")

	if owner == "" || repo == "" || name == "" {
		http.Error(w, "Missing owner, repo, or secret name", http.StatusBadRequest)
		return
	}

	userTs := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: user.AccessToken},
	)
	userTc := oauth2.NewClient(r.Context(), userTs)
	userGhClient := github.NewClient(userTc)
	hasAccess, err := app.repositories.HasMaintainerAccess(r.Context(), userGhClient, owner, repo)

	if err != nil {
		app.logger.Error("Failed to check permissions", slog.String("error", err.Error()))
		http.Error(w, "Permission check failed", http.StatusInternalServerError)
		return
	}
	if !hasAccess {
		app.logger.Warn("User attempted to delete secret without permission", slog.String("user", user.Email), slog.String("repo", owner+"/"+repo))
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Use Shared PAT Client (Only after verification)
	githubClient := app.patClient

	// Call repository service
	err = app.repositories.DeleteSecret(r.Context(), githubClient, owner, repo, name)
	if err != nil {
		app.logger.Error("Failed to delete secret", slog.String("error", err.Error()))
		http.Error(w, "Failed to delete secret", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (app *application) handleCreateSecret(w http.ResponseWriter, r *http.Request) {
	// Get Session & Verify User
	user, ok := app.requireUser(w, r)
	if !ok {
		return
	}

	// Extract path parameters
	owner := r.PathValue("owner")
	repo := r.PathValue("repo")
	// secret name from path
	name := r.PathValue("name")

	if owner == "" || repo == "" || name == "" {
		http.Error(w, "Missing owner, repo, or secret name", http.StatusBadRequest)
		return
	}

	// Parse Body
	var req struct {
		Value string `json:"value"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if req.Value == "" {
		http.Error(w, "Secret value is required", http.StatusBadRequest)
		return
	}

	// We do not work on userGhClient directly, but instead pass it to the repository service.
	userGhClient, err := app.getGitHubClient(r.Context(), user.AccessToken)
	if err != nil {
		app.logger.Error("Failed to create GitHub client", slog.String("error", err.Error()))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	hasAccess, err := app.repositories.HasMaintainerAccess(r.Context(), userGhClient, owner, repo)
	if err != nil {
		app.logger.Error("Failed to check permissions", slog.String("error", err.Error()))
		http.Error(w, "Permission check failed", http.StatusInternalServerError)
		return
	}
	if !hasAccess {
		app.logger.Warn("User attempted to create secret without permission", slog.String("user", user.Email), slog.String("repo", owner+"/"+repo))
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Use Shared PAT Client (Only after verification)
	githubClient := app.patClient

	err = app.repositories.CreateOrUpdateSecret(r.Context(), githubClient, owner, repo, name, req.Value)
	if err != nil {
		app.logger.Error("Failed to create secret", slog.String("error", err.Error()))
		http.Error(w, "Failed to create secret", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
