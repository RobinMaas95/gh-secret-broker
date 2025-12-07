package repository

import (
	"context"
	"strings"

	"github.com/google/go-github/v66/github"
)

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
				if repo.Permissions != nil && (repo.Permissions["admin"] || repo.Permissions["maintain"]) {
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
