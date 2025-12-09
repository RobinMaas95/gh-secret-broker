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

func TestHandleListSecrets(t *testing.T) {
	t.Run("Authorized and Has Access", func(t *testing.T) {
		// Mock Service
		mockService := &mockRepositoryService{
			HasMaintainerAccessFunc: func(ctx context.Context, client *github.Client, owner, repo string) (bool, error) {
				return true, nil
			},
			ListSecretsFunc: func(ctx context.Context, client *github.Client, owner, repo string) ([]string, error) {
				return []string{"SECRET_1", "SECRET_2"}, nil
			},
		}

		app := &application{
			logger:       setupTestLogger(),
			repositories: mockService,
			config:       &config.Config{GithubOrg: "test-org"},
		}

		// Mock Session
		store := sessions.NewCookieStore([]byte("secret"))
		gothic.Store = store
		user := goth.User{AccessToken: "valid-token"}

		req, _ := http.NewRequest("GET", "/api/repo/TargetOrg/repo-1/secrets", nil)
		// We need to set the path values (used by the standard library mux) manually
		// because normally the router would parse the url an populate this values
		// but in our tests no router is used
		req.SetPathValue("owner", "TargetOrg")
		req.SetPathValue("repo", "repo-1")

		w := httptest.NewRecorder()

		// Inject session
		session, _ := store.Get(req, "session")
		session.Values["user"] = user
		_ = session.Save(req, w)

		req.Header.Set("Cookie", w.Header().Get("Set-Cookie"))

		app.handleListSecrets(w, req)

		res := w.Result()
		defer func() { _ = res.Body.Close() }()

		assert.Equal(t, res.StatusCode, http.StatusOK)

		var secrets []string
		_ = json.NewDecoder(res.Body).Decode(&secrets)
		if len(secrets) != 2 {
			t.Errorf("expected 2 secrets, got %d", len(secrets))
		}
	})

	t.Run("Access Denied", func(t *testing.T) {
		// Mock Service returns false for access
		mockService := &mockRepositoryService{
			HasMaintainerAccessFunc: func(ctx context.Context, client *github.Client, owner, repo string) (bool, error) {
				return false, nil
			},
		}

		app := &application{
			logger:       setupTestLogger(),
			repositories: mockService,
			config:       &config.Config{GithubOrg: "test-org"},
		}

		store := sessions.NewCookieStore([]byte("secret"))
		gothic.Store = store
		user := goth.User{AccessToken: "valid-token"}

		req, _ := http.NewRequest("GET", "/api/repo/TargetOrg/repo-1/secrets", nil)
		req.SetPathValue("owner", "TargetOrg")
		req.SetPathValue("repo", "repo-1")
		w := httptest.NewRecorder()

		session, _ := store.Get(req, "session")
		session.Values["user"] = user
		_ = session.Save(req, w)
		req.Header.Set("Cookie", w.Header().Get("Set-Cookie"))

		app.handleListSecrets(w, req)

		res := w.Result()
		defer func() { _ = res.Body.Close() }()

		assert.Equal(t, res.StatusCode, http.StatusForbidden)
	})
}

func TestHandleDeleteSecret(t *testing.T) {
	t.Run("Authorized Delete", func(t *testing.T) {
		mockService := &mockRepositoryService{
			HasMaintainerAccessFunc: func(ctx context.Context, client *github.Client, owner, repo string) (bool, error) {
				return true, nil
			},
			DeleteSecretFunc: func(ctx context.Context, client *github.Client, owner, repo, name string) error {
				if owner == "TargetOrg" && repo == "repo-1" && name == "SECRET_1" {
					return nil
				}
				panic("unexpected arguments to DeleteSecret")
			},
		}

		app := &application{
			logger:       setupTestLogger(),
			repositories: mockService,
			config:       &config.Config{GithubOrg: "test-org"},
		}

		store := sessions.NewCookieStore([]byte("secret"))
		gothic.Store = store
		user := goth.User{AccessToken: "valid-token"}

		req, _ := http.NewRequest("DELETE", "/api/repo/TargetOrg/repo-1/secrets/SECRET_1", nil)
		req.SetPathValue("owner", "TargetOrg")
		req.SetPathValue("repo", "repo-1")
		req.SetPathValue("name", "SECRET_1")
		w := httptest.NewRecorder()

		session, _ := store.Get(req, "session")
		session.Values["user"] = user
		_ = session.Save(req, w)

		req.Header.Set("Cookie", w.Header().Get("Set-Cookie"))

		app.handleDeleteSecret(w, req)

		res := w.Result()
		defer func() { _ = res.Body.Close() }()

		assert.Equal(t, res.StatusCode, http.StatusNoContent)
	})
}
