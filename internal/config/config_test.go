package config

import (
	"os"
	"strings"
	"testing"
)

func TestConfig_Load(t *testing.T) {
	tests := []struct {
		name        string
		envs        map[string]string
		wantErr     bool
		errContains string
	}{
		{
			name: "All required envs set",
			envs: map[string]string{
				"GITHUB_CLIENT_ID":     "test-client-id",
				"GITHUB_CLIENT_SECRET": "test-client-secret",
			},
			wantErr: false,
		},
		{
			name: "With session secret and production mode",
			envs: map[string]string{
				"GITHUB_CLIENT_ID":     "test-client-id",
				"GITHUB_CLIENT_SECRET": "test-client-secret",
				"SESSION_SECRET":       "this-is-a-very-long-session-secret-key",
				"ENVIRONMENT":          "production",
			},
			wantErr: false,
		},
		{
			name: "Missing client ID",
			envs: map[string]string{
				"GITHUB_CLIENT_SECRET": "test-client-secret",
			},
			wantErr:     true,
			errContains: "GITHUB_CLIENT_ID",
		},
		{
			name: "Missing client secret",
			envs: map[string]string{
				"GITHUB_CLIENT_ID": "test-client-id",
			},
			wantErr:     true,
			errContains: "GITHUB_CLIENT_SECRET",
		},
		{
			name:        "Missing all required envs",
			envs:        map[string]string{},
			wantErr:     true,
			errContains: "GITHUB_CLIENT_ID",
		},
		{
			name: "Production without session secret",
			envs: map[string]string{
				"GITHUB_CLIENT_ID":     "test-client-id",
				"GITHUB_CLIENT_SECRET": "test-client-secret",
				"ENVIRONMENT":          "production",
			},
			wantErr:     true,
			errContains: "SESSION_SECRET is required",
		},
		{
			name: "Session secret too short",
			envs: map[string]string{
				"GITHUB_CLIENT_ID":     "test-client-id",
				"GITHUB_CLIENT_SECRET": "test-client-secret",
				"SESSION_SECRET":       "short",
			},
			wantErr:     true,
			errContains: "at least 32 characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save and clear all relevant environment variables
			envKeys := []string{
				"GITHUB_CLIENT_ID",
				"GITHUB_CLIENT_SECRET",
				"SESSION_SECRET",
				"ENVIRONMENT",
				"BASE_URL",
			}
			originalEnvs := make(map[string]string)
			for _, key := range envKeys {
				originalEnvs[key] = os.Getenv(key)
				_ = os.Unsetenv(key) //nolint:errcheck
			}

			// Restore environment after test
			t.Cleanup(func() {
				for key, val := range originalEnvs {
					if val != "" {
						_ = os.Setenv(key, val) //nolint:errcheck
					} else {
						_ = os.Unsetenv(key) //nolint:errcheck
					}
				}
			})

			// Set test environment variables
			for key, val := range tt.envs {
				_ = os.Setenv(key, val) //nolint:errcheck
			}

			got, err := Load()

			if tt.wantErr {
				if err == nil {
					t.Errorf("Load() expected error containing %q, got nil", tt.errContains)
					return
				}
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("Load() error = %v, want error containing %q", err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Fatalf("Load() unexpected error = %v", err)
			}

			// Verify expected values
			if got.GithubClientID != tt.envs["GITHUB_CLIENT_ID"] {
				t.Errorf("GithubClientID = %v, want %v", got.GithubClientID, tt.envs["GITHUB_CLIENT_ID"])
			}
			if got.GithubClientSecret != tt.envs["GITHUB_CLIENT_SECRET"] {
				t.Errorf("GithubClientSecret = %v, want %v", got.GithubClientSecret, tt.envs["GITHUB_CLIENT_SECRET"])
			}

			// Check session secret is set (either from env or generated)
			if got.SessionSecret == "" {
				t.Error("SessionSecret should not be empty")
			}

			// Check environment defaults to development
			expectedEnv := tt.envs["ENVIRONMENT"]
			if expectedEnv == "" {
				expectedEnv = "development"
			}
			if got.Environment != expectedEnv {
				t.Errorf("Environment = %v, want %v", got.Environment, expectedEnv)
			}
		})
	}
}

func TestConfig_IsProduction(t *testing.T) {
	tests := []struct {
		env  string
		want bool
	}{
		{"production", true},
		{"development", false},
		{"", false},
		{"staging", false},
	}

	for _, tt := range tests {
		t.Run(tt.env, func(t *testing.T) {
			cfg := &Config{Environment: tt.env}
			if got := cfg.IsProduction(); got != tt.want {
				t.Errorf("IsProduction() = %v, want %v", got, tt.want)
			}
		})
	}
}
