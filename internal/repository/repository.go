package repository

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/google/go-github/v80/github"
	"golang.org/x/crypto/nacl/box"
)

// RepositoryService defines the interface for repository operations.
// This allows for mocking in tests.
type RepositoryService interface {
	ListMaintainableRepositories(ctx context.Context, client *github.Client, orgName string) ([]*github.Repository, error)
	ListSecrets(ctx context.Context, client *github.Client, owner, repo string) ([]string, error)
	DeleteSecret(ctx context.Context, client *github.Client, owner, repo, name string) error
	CreateOrUpdateSecret(ctx context.Context, client *github.Client, owner, repo, name, value string) error
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

// DeleteSecret deletes a secret from a repository.
func (s *Service) DeleteSecret(ctx context.Context, client *github.Client, owner, repo, name string) error {
	_, err := client.Actions.DeleteRepoSecret(ctx, owner, repo, name)
	return err
}

// CreateOrUpdateSecret encrypts and uploads a secret to a repository.
func (s *Service) CreateOrUpdateSecret(ctx context.Context, client *github.Client, owner, repo, name, value string) error {
	// 1. Get Public Key from GitHub
	publicKey, _, err := client.Actions.GetRepoPublicKey(ctx, owner, repo)
	if err != nil {
		return err
	}

	// 2. Encrypt the secret
	encryptedValue, err := encryptSecretWithPublicKey(publicKey, name, value)
	if err != nil {
		return err
	}

	// 3. Create or Update Secret
	secret := &github.EncryptedSecret{
		Name:           name,
		KeyID:          publicKey.GetKeyID(),
		EncryptedValue: encryptedValue,
	}

	_, err = client.Actions.CreateOrUpdateRepoSecret(ctx, owner, repo, secret)
	return err
}

// encryptSecretWithPublicKey encrypts a secret value using the repository's public key (NaCl Box).
// Note: This logic mimics the standard GitHub Actions encryption process.
func encryptSecretWithPublicKey(publicKey *github.PublicKey, secretName, secretValue string) (string, error) {
	if publicKey == nil || publicKey.Key == nil {
		return "", fmt.Errorf("public key is missing")
	}

	// Decode the base64 public key
	decodedKey, err := base64.StdEncoding.DecodeString(*publicKey.Key)
	if err != nil {
		return "", fmt.Errorf("failed to decode public key: %w", err)
	}

	// Convert fixed size key
	var recipientKey [32]byte
	copy(recipientKey[:], decodedKey)

	// Encrypt using NaCl Box (Sealed Box)
	// GitHub uses libsodium's crypto_box_seal.
	// In Go's x/crypto/nacl/box, we use box.SealAnonymous
	encryptedBytes, err := box.SealAnonymous(nil, []byte(secretValue), &recipientKey, rand.Reader)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt secret: %w", err)
	}

	// Encode result to base64
	return base64.StdEncoding.EncodeToString(encryptedBytes), nil
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
