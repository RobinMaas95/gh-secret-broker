package config

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestConfig_Load(t *testing.T) {
	tests := []struct {
		name    string // description of this test case
		setEnvs bool
		want    *Config
	}{
		{
			name:    "Envs found",
			setEnvs: true,
			want: &Config{
				GithubClientID:     "DUMMY_ID",
				GithubClientSecret: "DUMMY_SECRET",
			},
		},
		{
			name:    "Envs not found",
			setEnvs: false,
			want: &Config{
				GithubClientID:     "",
				GithubClientSecret: "",
			},
		},
	}
	for _, tt := range tests {
		// Save original value of envs
		originalGITHUB_CLIENT_ID := os.Getenv("GITHUB_CLIENT_ID")
		originalGITHUB_CLIENT_SECRET := os.Getenv("GITHUB_CLIENT_SECRET")

		t.Run(tt.name, func(t *testing.T) {
			if tt.setEnvs {
				err := os.Setenv("GITHUB_CLIENT_ID", "DUMMY_ID")
				if err != nil {
					t.Fatal(err)
				}

				err = os.Setenv("GITHUB_CLIENT_SECRET", "DUMMY_SECRET")
				if err != nil {
					t.Fatal(err)
				}

			}
			t.Cleanup(func() {
				err := os.Setenv("GITHUB_CLIENT_ID", originalGITHUB_CLIENT_ID)
				if err != nil {
					t.Fatal(err)
				}

				err = os.Setenv("GITHUB_CLIENT_SECRET", originalGITHUB_CLIENT_SECRET)
				if err != nil {
					t.Fatal(err)
				}
			})
			got, err := Load()
			if err != nil {
				if tt.name != "Envs not found" {
					t.Fatalf("Load() error = %v, want nil", err)
				}
				return
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("Config mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
