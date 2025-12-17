package main

import (
	"context"

	"github.com/google/go-github/v80/github"
)

// mockRepositoryService mocks the repository.RepositoryService interface
type mockRepositoryService struct {
	ListMaintainableRepositoriesFunc func(ctx context.Context, client *github.Client, orgName string) ([]*github.Repository, error)
	ListSecretsFunc                  func(ctx context.Context, client *github.Client, owner, repo string) ([]string, error)
	DeleteSecretFunc                 func(ctx context.Context, client *github.Client, owner, repo, name string) error
	CreateOrUpdateSecretFunc         func(ctx context.Context, client *github.Client, owner, repo, name, value string) error
	HasMaintainerAccessFunc          func(ctx context.Context, client *github.Client, owner, repo string) (bool, error)
}

// Our interfaces only check if a function is set, if so they call it, otherwise they return nil.
func (m *mockRepositoryService) ListMaintainableRepositories(ctx context.Context, client *github.Client, orgName string) ([]*github.Repository, error) {
	if m.ListMaintainableRepositoriesFunc != nil {
		return m.ListMaintainableRepositoriesFunc(ctx, client, orgName)
	}
	return nil, nil
}

func (m *mockRepositoryService) ListSecrets(ctx context.Context, client *github.Client, owner, repo string) ([]string, error) {
	if m.ListSecretsFunc != nil {
		return m.ListSecretsFunc(ctx, client, owner, repo)
	}
	return nil, nil
}

func (m *mockRepositoryService) DeleteSecret(ctx context.Context, client *github.Client, owner, repo, name string) error {
	if m.DeleteSecretFunc != nil {
		return m.DeleteSecretFunc(ctx, client, owner, repo, name)
	}
	return nil
}

func (m *mockRepositoryService) CreateOrUpdateSecret(ctx context.Context, client *github.Client, owner, repo, name, value string) error {
	if m.CreateOrUpdateSecretFunc != nil {
		return m.CreateOrUpdateSecretFunc(ctx, client, owner, repo, name, value)
	}
	return nil
}

func (m *mockRepositoryService) HasMaintainerAccess(ctx context.Context, client *github.Client, owner, repo string) (bool, error) {
	if m.HasMaintainerAccessFunc != nil {
		return m.HasMaintainerAccessFunc(ctx, client, owner, repo)
	}
	return false, nil
}
