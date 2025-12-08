package config

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	Addr               string
	BaseURL            string
	Environment        string // "development" or "production"
	SessionSecret      string
	GithubClientID     string
	GithubClientSecret string
	GithubOrg          string
	GithubPAT          string
}

// IsProduction returns true if running in production environment
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

func Load() (*Config, error) {
	var errs []error
	getEnv := func(key string) string {
		val := os.Getenv(key)
		if val == "" {
			errs = append(errs, fmt.Errorf("environment variable %s is not set", key))
		}
		return val
	}

	// Get environment, default to development
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = "development"
	}

	// Get session secret - required in production, generated in development
	sessionSecret := os.Getenv("SESSION_SECRET")
	if sessionSecret == "" {
		if env == "production" {
			errs = append(errs, fmt.Errorf("SESSION_SECRET is required in production"))
		} else {
			// Generate a random secret for development
			sessionSecret = generateRandomSecret()
		}
	} else if len(sessionSecret) < 32 {
		errs = append(errs, fmt.Errorf("SESSION_SECRET must be at least 32 characters"))
	}

	config := &Config{
		BaseURL:            os.Getenv("BASE_URL"), // Optional
		Environment:        env,
		SessionSecret:      sessionSecret,
		GithubClientID:     getEnv("GITHUB_CLIENT_ID"),
		GithubClientSecret: getEnv("GITHUB_CLIENT_SECRET"),
		GithubOrg:          getEnv("GITHUB_ORG"),
		GithubPAT:          getEnv("GITHUB_PAT"),
	}

	if len(errs) > 0 {
		errMsgs := make([]string, len(errs))
		for i, e := range errs {
			errMsgs[i] = e.Error()
		}
		return nil, errors.New(strings.Join(errMsgs, "; "))
	}

	return config, nil
}

// generateRandomSecret generates a cryptographically secure random secret
func generateRandomSecret() string {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		panic("failed to generate random secret: " + err.Error())
	}
	return base64.StdEncoding.EncodeToString(bytes)
}
