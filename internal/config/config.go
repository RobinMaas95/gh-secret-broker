package config

import (
	"fmt"
	"os"
)

type Config struct {
	Addr               string
	GithubClientID     string
	GithubClientSecret string
	// staticDir string
}

func Load() (*Config, error) {
	var err error
	getEnv := func(key string) string {
		val := os.Getenv(key)
		if val == "" {
			err = fmt.Errorf("environment variable %s is not set", key)
		}
		return val
	}

	config := &Config{
		GithubClientID:     getEnv("GITHUB_CLIENT_ID"),
		GithubClientSecret: getEnv("GITHUB_CLIENT_SECRET"),
	}

	if err != nil {
		return nil, err
	}

	return config, nil
}
