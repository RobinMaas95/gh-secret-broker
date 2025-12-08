package repository

import (
	"context"
	"strings"

	"github.com/google/go-github/v80/github"
)

// RepositoryService defines the interface for repository operations.
// This allows for mocking in tests.
type RepositoryService interface {
	ListMaintainableRepositories(ctx context.Context, client *github.Client, orgName string) ([]*github.Repository, error)
	ListSecrets(ctx context.Context, client *github.Client, owner, repo string) ([]string, error)
	HasMaintainerAccess(ctx context.Context, client *github.Client, owner, repo string) (bool, error)
}

type Service struct{}

func NewService() *Service {
	return &Service{}
}

// ListMaintainableRepositories lists all repositories in the given organization
// where the authenticated user has 'maintain' or 'admin' permissions.
func (s *Service) ListMaintainableRepositories(ctx context.Context, client *github.Client, orgName string) ([]*github.Repository, error) {
	opts := &github.RepositoryListByAuthenticatedUserOptions{
		Type:        "all",
		ListOptions: github.ListOptions{PerPage: 100},
	}

	var allRepos []*github.Repository
	for {
		repos, resp, err := client.Repositories.ListByAuthenticatedUser(ctx, opts)
		if err != nil {
			return nil, err
		}

		for _, repo := range repos {
			// Filter by Organization
			if repo.Owner != nil && repo.Owner.Login != nil && strings.EqualFold(*repo.Owner.Login, orgName) {
				// Filter by Permissions (Admin or Maintain)
				// Note: Maintain permission is often implied by Admin, but explicit check is safer.
				if s.hasMaintainerPermissions(repo) {
					allRepos = append(allRepos, repo)
				}
			}
		}

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return allRepos, nil
}

// ListSecrets lists the names of secrets for a repository.
// Note: GitHub API does not return secret values, only names and metadata.
func (s *Service) ListSecrets(ctx context.Context, client *github.Client, owner, repo string) ([]string, error) {
	opts := &github.ListOptions{PerPage: 100}
	allSecrets := []string{}

	for {
		secrets, resp, err := client.Actions.ListRepoSecrets(ctx, owner, repo, opts)
		if err != nil {
			return nil, err
		}

		for _, secret := range secrets.Secrets {
			allSecrets = append(allSecrets, secret.Name)
		}

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return allSecrets, nil
}

// HasMaintainerAccess checks if the user has 'admin' or 'maintain' permissions on the repository.
func (s *Service) HasMaintainerAccess(ctx context.Context, client *github.Client, owner, repo string) (bool, error) {
	repository, _, err := client.Repositories.Get(ctx, owner, repo)
	if err != nil {
		return false, err
	}

	if repository.Permissions == nil {
		return false, nil
	}

	return s.hasMaintainerPermissions(repository), nil
}

func (s *Service) hasMaintainerPermissions(repo *github.Repository) bool {
	permissions := repo.GetPermissions()
	return permissions["admin"] || permissions["maintain"]
}
