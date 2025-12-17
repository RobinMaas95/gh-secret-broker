package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// handleListRepositories handles the GET /api/repositories request.
// It retrieves the list of repositories where the user has maintain/admin access
// in the organization configured in GITHUB_ORG.
func (app *application) handleListRepositories(w http.ResponseWriter, r *http.Request) {
	// Get Session & Verify User
	user, ok := app.requireUser(w, r)
	if !ok {
		return
	}

	// Create GitHub Client using User's Token
	githubClient, err := app.getGitHubClient(r.Context(), user.AccessToken)
	if err != nil {
		app.logger.Error("Failed to create GitHub client", slog.String("error", err.Error()))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Retrieve Repositories via Service
	orgName := app.config.GithubOrg
	if orgName == "" {
		app.logger.Error("GITHUB_ORG is not configured")
		http.Error(w, "Configuration Error: GITHUB_ORG not set", http.StatusInternalServerError)
		return
	}

	repos, err := app.repositories.ListMaintainableRepositories(r.Context(), githubClient, orgName)
	if err != nil {
		app.logger.Error("Failed to list repositories", slog.String("error", err.Error()), slog.String("org", orgName))
		http.Error(w, "Failed to fetch repositories", http.StatusInternalServerError)
		return
	}

	// Respond with JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(repos); err != nil {
		app.logger.Error("Failed to encode response", slog.String("error", err.Error()))
	}
}
