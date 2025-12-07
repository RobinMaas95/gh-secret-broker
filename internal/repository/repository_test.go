package repository_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RobinMaas95/gh-secret-broker/internal/repository"
	"github.com/google/go-github/v66/github"
)

func TestListUserRepositories(t *testing.T) {
	// Setup mock server
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	// Define expected behavior
	mux.HandleFunc("/user/repos", func(w http.ResponseWriter, r *http.Request) {
		// Verify parameters
		if r.URL.Query().Get("type") != "all" {
			t.Errorf("expected type=all, got %s", r.URL.Query().Get("type"))
		}

		// Mock response
		repos := []*github.Repository{
			{
				Name: github.String("repo-1"),
				Owner: &github.User{
					Login: github.String("TargetOrg"),
				},
				Permissions: map[string]bool{
					"admin":    true,
					"maintain": true,
				},
			},
			{
				Name: github.String("repo-2"),
				Owner: &github.User{
					Login: github.String("OtherOrg"),
				},
				Permissions: map[string]bool{
					"admin": true,
				},
			},
			{
				Name: github.String("repo-3"),
				Owner: &github.User{
					Login: github.String("TargetOrg"),
				},
				Permissions: map[string]bool{
					"admin":    false,
					"maintain": false,
					"push":     true,
				},
			},
		}
		_ = json.NewEncoder(w).Encode(repos)
	})

	// Setup client to use the mock server
	client := github.NewClient(nil)
	client.BaseURL, _ = client.BaseURL.Parse(server.URL + "/")

	// This part is a bit tricky because usually we pass a token source to NewClient,
	// but here we just want to test the logic given a client.
	// So our service should probably accept the client?
	// Or we configure the service to use a specific base URL for testing?

	service := repository.NewService()

	// We'll pass the mocked client directly for this test helper if possible,
	// or we mock the "Token to Client" creation.
	// For simplicity, let's assume the method takes a *github.Client.
	repos, err := service.ListMaintainableRepositories(context.Background(), client, "TargetOrg")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(repos) != 1 {
		t.Errorf("expected 1 repo, got %d", len(repos))
	}

	if repos[0].GetName() != "repo-1" {
		t.Errorf("expected repo-1, got %s", repos[0].GetName())
	}
}
