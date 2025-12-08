package repository_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RobinMaas95/gh-secret-broker/internal/repository"
	"github.com/google/go-github/v80/github"
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
				Name: github.Ptr("repo-1"),
				Owner: &github.User{
					Login: github.Ptr("TargetOrg"),
				},
				Permissions: map[string]bool{
					"admin":    true,
					"maintain": true,
				},
			},
			{
				Name: github.Ptr("repo-2"),
				Owner: &github.User{
					Login: github.Ptr("OtherOrg"),
				},
				Permissions: map[string]bool{
					"admin": true,
				},
			},
			{
				Name: github.Ptr("repo-3"),
				Owner: &github.User{
					Login: github.Ptr("TargetOrg"),
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

func TestListSecrets(t *testing.T) {
	// Setup mock server
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	// Define fake GitHub server, we will call this in our test
	mux.HandleFunc("/repos/TargetOrg/repo-1/actions/secrets", func(w http.ResponseWriter, r *http.Request) {
		response := &github.Secrets{
			TotalCount: 2,
			Secrets: []*github.Secret{
				{Name: "SECRET_ONE"},
				{Name: "SECRET_TWO"},
			},
		}
		_ = json.NewEncoder(w).Encode(response)
	})

	// Setup client
	client := github.NewClient(nil)
	// Overwriting target (GitHub server) url, so that our fake server is used
	client.BaseURL, _ = client.BaseURL.Parse(server.URL + "/")

	service := repository.NewService()
	secrets, err := service.ListSecrets(context.Background(), client, "TargetOrg", "repo-1")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(secrets) != 2 {
		t.Errorf("expected 2 secrets, got %d", len(secrets))
	}

	expected := []string{"SECRET_ONE", "SECRET_TWO"}
	for i, s := range secrets {
		if s != expected[i] {
			t.Errorf("expected secret %s, got %s", expected[i], s)
		}
	}
}

func TestHasMaintainerAccess(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	mux.HandleFunc("/repos/TargetOrg/repo-1", func(w http.ResponseWriter, r *http.Request) {
		repo := &github.Repository{
			Permissions: map[string]bool{"admin": true},
		}
		_ = json.NewEncoder(w).Encode(repo)
	})

	mux.HandleFunc("/repos/TargetOrg/repo-2", func(w http.ResponseWriter, r *http.Request) {
		repo := &github.Repository{
			Permissions: map[string]bool{"pull": true}, // No admin/maintain
		}
		_ = json.NewEncoder(w).Encode(repo)
	})

	client := github.NewClient(nil)
	client.BaseURL, _ = client.BaseURL.Parse(server.URL + "/")
	service := repository.NewService()

	// Case 1: Has Access
	hasAccess, err := service.HasMaintainerAccess(context.Background(), client, "TargetOrg", "repo-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !hasAccess {
		t.Error("expected access, got denied")
	}

	// Case 2: No Access
	hasAccess, err = service.HasMaintainerAccess(context.Background(), client, "TargetOrg", "repo-2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if hasAccess {
		t.Error("expected denied, got access")
	}
}
