package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RobinMaas95/gh-secret-broker/internal/assert"
	"github.com/RobinMaas95/gh-secret-broker/internal/config"
	"github.com/google/go-github/v80/github"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
)

func TestHandleListRepositories(t *testing.T) {
	t.Run("Authorized user", func(t *testing.T) {
		// Mock Service
		mockService := &mockRepositoryService{
			ListMaintainableRepositoriesFunc: func(ctx context.Context, client *github.Client, orgName string) ([]*github.Repository, error) {
				return []*github.Repository{
					{Name: github.Ptr("repo-1")},
				}, nil
			},
		}

		app := &application{
			logger:       setupTestLogger(),
			repositories: mockService,
			config:       &config.Config{GithubOrg: "test-org"},
		}

		// Mock Session
		user := goth.User{AccessToken: "valid-token"}

		// Mock Session Store
		store := sessions.NewCookieStore([]byte("secret"))
		gothic.Store = store

		// Create a request
		req, _ := http.NewRequest("GET", "/api/repositories", nil)
		w := httptest.NewRecorder()

		// Inject session
		session, _ := store.Get(req, "session")
		session.Values["user"] = user
		_ = session.Save(req, w)

		// Create a new request with the cookie
		req, _ = http.NewRequest("GET", "/api/repositories", nil)
		req.Header.Set("Cookie", w.Header().Get("Set-Cookie"))

		// Run handler
		app.handleListRepositories(w, req)

		res := w.Result()
		defer func() { _ = res.Body.Close() }()

		assert.Equal(t, res.StatusCode, http.StatusOK)

		var repos []*github.Repository
		_ = json.NewDecoder(res.Body).Decode(&repos)
		if len(repos) != 1 {
			t.Errorf("expected 1 repo, got %d", len(repos))
		}
		if *repos[0].Name != "repo-1" {
			t.Errorf("expected repo-1, got %s", *repos[0].Name)
		}
	})

	t.Run("Unauthorized user", func(t *testing.T) {
		app := &application{
			logger:       setupTestLogger(),
			repositories: &mockRepositoryService{},
			config:       &config.Config{GithubOrg: "test-org"},
		}

		req, _ := http.NewRequest("GET", "/api/repositories", nil)
		w := httptest.NewRecorder()

		app.handleListRepositories(w, req)
		res := w.Result()
		defer func() { _ = res.Body.Close() }()

		assert.Equal(t, res.StatusCode, http.StatusUnauthorized)
	})
}
