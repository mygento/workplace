package runtime

import (
	"testing"

	"github.com/mygento/workplace/internal/config"
)

func TestDockerPHPContainerName(t *testing.T) {
	tests := []struct {
		name        string
		cfg         config.Config
		expected    string
		expectError bool
	}{
		{
			name: "default container name",
			cfg: config.Config{
				Name: "my-project",
				Services: map[string]config.Service{
					"php": {},
				},
			},
			expected: "my-project-php-1",
		},
		{
			name: "custom container name",
			cfg: config.Config{
				Name: "my-project",
				Services: map[string]config.Service{
					"php": {
						ContainerName: "custom-php",
					},
				},
			},
			expected: "custom-php",
		},
		{
			name: "missing PHP service",
			cfg: config.Config{
				Name:     "my-project",
				Services: map[string]config.Service{},
			},
			expectError: true,
		},
		{
			name: "missing project name",
			cfg: config.Config{
				Services: map[string]config.Service{
					"php": {},
				},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := dockerPHPContainerName(tt.cfg)

			if tt.expectError {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if actual != tt.expected {
				t.Errorf(
					"expected %q, got %q",
					tt.expected,
					actual,
				)
			}
		})
	}
}
