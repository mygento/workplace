package environment

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestMigrateConfig(t *testing.T) {
	tests := []struct {
		name            string
		packageJSON     string
		existing        string
		localJSON       string
		existingLocal   string
		overwrite       bool
		wantErr         error
		wantType        string
		wantWritten     bool
		wantConfig      map[string]interface{}
		wantLocalConfig map[string]interface{}
	}{
		{
			name: "migrates basic configuration",
			packageJSON: `{
				"name": "test-project",
				"workplace": {
					"type": "magento2",
					"php": {
						"image": "php:8.2"
					},
					"nginx": {
						"port": 8080
					}
				}
			}`,
			wantType:    "magento2",
			wantWritten: true,
			wantConfig: map[string]interface{}{
				"name": "test-project",
				"type": "magento2",
				"services": map[string]interface{}{
					"php": map[string]interface{}{
						"image": "php:8.2",
					},
					"nginx": map[string]interface{}{
						"port": 8080,
					},
				},
			},
		},
		{
			name: "migrates local configuration",
			packageJSON: `{
				"name": "test-project",
				"workplace": {
					"type": "magento2",
					"php": {
						"image": "php:8.2"
					}
				}
			}`,
			localJSON: `{
				"php": {
					"image": "ghcr.io/mygento/php:8.1-debug"
				},
				"nginx": {
					"port": 8082
				}
			}`,
			wantType:    "magento2",
			wantWritten: true,
			wantConfig: map[string]interface{}{
				"name": "test-project",
				"type": "magento2",
				"services": map[string]interface{}{
					"php": map[string]interface{}{
						"image": "php:8.2",
					},
				},
			},
			wantLocalConfig: map[string]interface{}{
				"services": map[string]interface{}{
					"php": map[string]interface{}{
						"image": "ghcr.io/mygento/php:8.1-debug",
					},
					"nginx": map[string]interface{}{
						"port": 8082,
					},
				},
			},
		},
		{
			name: "uses directory name when project name is missing",
			packageJSON: `{
				"workplace": {
					"type": "magento2"
				}
			}`,
			wantType:    "magento2",
			wantWritten: true,
		},
		{
			name: "converts elasticsearch to elastic",
			packageJSON: `{
				"name": "test-project",
				"workplace": {
					"type": "magento2",
					"elasticsearch": {
						"image": "opensearchproject/opensearch:3"
					}
				}
			}`,
			wantType:    "magento2",
			wantWritten: true,
			wantConfig: map[string]interface{}{
				"name": "test-project",
				"type": "magento2",
				"services": map[string]interface{}{
					"elastic": map[string]interface{}{
						"image": "opensearchproject/opensearch:3",
					},
				},
			},
		},
		{
			name: "existing project destination",
			packageJSON: `{
				"name": "test-project",
				"workplace": {
					"type": "magento2"
				}
			}`,
			existing: "existing configuration",
			wantErr:  ErrTargetConfigExists,
		},
		{
			name: "overwrites existing project destination",
			packageJSON: `{
				"name": "test-project",
				"workplace": {
					"type": "magento2"
				}
			}`,
			existing:    "existing configuration",
			overwrite:   true,
			wantType:    "magento2",
			wantWritten: true,
		},
		{
			name: "local destination already exists",
			packageJSON: `{
				"workplace": {
					"type": "magento2"
				}
			}`,
			localJSON: `{
				"php": {
					"image": "php:8.1"
				}
			}`,
			existingLocal: "services: {}\n",
			wantErr:       ErrLocalConfigExists,
		},
		{
			name: "overwrites existing local destination",
			packageJSON: `{
				"workplace": {
					"type": "magento2"
				}
			}`,
			localJSON: `{
				"php": {
					"image": "ghcr.io/mygento/php:8.1-debug"
				}
			}`,
			existingLocal: "services: {}\n",
			overwrite:     true,
			wantType:      "magento2",
			wantWritten:   true,
			wantLocalConfig: map[string]interface{}{
				"services": map[string]interface{}{
					"php": map[string]interface{}{
						"image": "ghcr.io/mygento/php:8.1-debug",
					},
				},
			},
		},
		{
			name: "invalid local JSON",
			packageJSON: `{
				"workplace": {
					"type": "magento2"
				}
			}`,
			localJSON: `{invalid json}`,
		},
		{
			name:        "invalid package JSON",
			packageJSON: `{invalid json}`,
		},
		{
			name:        "missing workplace configuration",
			packageJSON: `{"name":"test-project"}`,
			wantErr:     ErrLegacyConfigNotFound,
		},
		{
			name: "missing project type",
			packageJSON: `{
				"name": "test-project",
				"workplace": {}
			}`,
		},
		{
			name: "migrates without local configuration",
			packageJSON: `{
				"name": "test-project",
				"workplace": {
					"type": "magento2"
				}
			}`,
			wantType:    "magento2",
			wantWritten: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()

			packageJSONPath := filepath.Join(dir, "package.json")

			if err := os.WriteFile(
				packageJSONPath,
				[]byte(tt.packageJSON),
				0644,
			); err != nil {
				t.Fatal(err)
			}

			workplaceDir := filepath.Join(dir, ".workplace")

			if tt.localJSON != "" || tt.existingLocal != "" {
				if err := os.Mkdir(workplaceDir, 0755); err != nil {
					t.Fatal(err)
				}
			}

			if tt.localJSON != "" {
				if err := os.WriteFile(
					filepath.Join(workplaceDir, "config.local.json"),
					[]byte(tt.localJSON),
					0644,
				); err != nil {
					t.Fatal(err)
				}
			}

			if tt.existingLocal != "" {
				if err := os.WriteFile(
					filepath.Join(workplaceDir, "config.local.yaml"),
					[]byte(tt.existingLocal),
					0644,
				); err != nil {
					t.Fatal(err)
				}
			}

			destination := filepath.Join(dir, ".workplace.yaml")

			if tt.existing != "" {
				if err := os.WriteFile(
					destination,
					[]byte(tt.existing),
					0644,
				); err != nil {
					t.Fatal(err)
				}
			}

			result, err := MigrateConfig(dir, tt.overwrite)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf(
						"expected error %v, got %v",
						tt.wantErr,
						err,
					)
				}
				return
			}

			if tt.name == "invalid local JSON" {
				if err == nil || !strings.Contains(
					err.Error(),
					"parse local config",
				) {
					t.Fatalf("expected local JSON error, got %v", err)
				}
				return
			}

			if tt.name == "invalid package JSON" {
				if err == nil {
					t.Fatal("expected package JSON parsing error")
				}
				return
			}

			if strings.Contains(tt.name, "missing project type") {
				if err == nil {
					t.Fatal("expected missing project type error")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected migration error: %v", err)
			}

			if result == nil {
				t.Fatal("expected migration result")
			}

			if result.Type != tt.wantType {
				t.Errorf(
					"unexpected type: got %q, want %q",
					result.Type,
					tt.wantType,
				)
			}

			if result.Written != tt.wantWritten {
				t.Errorf(
					"unexpected Written: got %v, want %v",
					result.Written,
					tt.wantWritten,
				)
			}

			// Validate project YAML.
			if tt.wantWritten {
				data, err := os.ReadFile(destination)
				if err != nil {
					t.Fatalf("read project YAML: %v", err)
				}

				var actual map[string]interface{}
				if err := yaml.Unmarshal(data, &actual); err != nil {
					t.Fatalf("parse project YAML: %v", err)
				}

				if tt.wantConfig != nil {
					assertYAMLEqual(t, tt.wantConfig, actual)
				}

				if tt.name == "uses directory name when project name is missing" {
					if actual["name"] != filepath.Base(dir) {
						t.Errorf(
							"unexpected project name: %v",
							actual["name"],
						)
					}
				}
			}

			// Validate local YAML.
			if tt.wantLocalConfig != nil {
				localDestination := filepath.Join(
					workplaceDir,
					"config.local.yaml",
				)

				data, err := os.ReadFile(localDestination)
				if err != nil {
					t.Fatalf("read local YAML: %v", err)
				}

				var actual map[string]interface{}
				if err := yaml.Unmarshal(data, &actual); err != nil {
					t.Fatalf("parse local YAML: %v", err)
				}

				assertYAMLEqual(t, tt.wantLocalConfig, actual)
			}

			// Original local JSON must never be deleted.
			if tt.localJSON != "" {
				if _, err := os.Stat(
					filepath.Join(workplaceDir, "config.local.json"),
				); err != nil {
					t.Errorf("legacy local JSON was removed: %v", err)
				}
			}

			// Check no temporary files remain.
			assertNoTemporaryFiles(t, dir)
		})
	}
}

func TestMigrateServices(t *testing.T) {
	input := map[string]json.RawMessage{
		"php":           json.RawMessage(`{"image":"php:8.3"}`),
		"nginx":         json.RawMessage(`{"port":8082}`),
		"type":          json.RawMessage(`"magento2"`),
		"elasticsearch": json.RawMessage(`{"image":"elastic:latest"}`),
		"empty":         json.RawMessage(`null`),
	}

	services, err := migrateServices(input, "type")
	if err != nil {
		t.Fatalf("migrate services: %v", err)
	}

	if _, ok := services["type"]; ok {
		t.Fatal("type should not be included")
	}

	if _, ok := services["empty"]; ok {
		t.Fatal("null service should not be included")
	}

	if _, ok := services["elasticsearch"]; ok {
		t.Fatal("legacy elasticsearch name should be converted")
	}

	if _, ok := services["elastic"]; !ok {
		t.Fatal("elastic service missing")
	}

	if _, ok := services["php"]; !ok {
		t.Fatal("PHP service missing")
	}

	if _, ok := services["nginx"]; !ok {
		t.Fatal("nginx service missing")
	}
}

func TestMigrateServicesInvalidJSON(t *testing.T) {
	input := map[string]json.RawMessage{
		"php": json.RawMessage(`{invalid}`),
	}

	_, err := migrateServices(input)
	if err == nil {
		t.Fatal("expected parsing error")
	}
}

func TestMigrateConfigPreservesLocalJSON(t *testing.T) {
	dir := t.TempDir()

	workplaceDir := filepath.Join(dir, ".workplace")
	if err := os.Mkdir(workplaceDir, 0755); err != nil {
		t.Fatal(err)
	}

	packageJSON := `{
		"name": "test-project",
		"workplace": {
			"type": "magento2"
		}
	}`

	if err := os.WriteFile(
		filepath.Join(dir, "package.json"),
		[]byte(packageJSON),
		0644,
	); err != nil {
		t.Fatal(err)
	}

	localJSON := `{
		"php": {
			"image": "ghcr.io/mygento/php:8.1-debug"
		}
	}`

	localSource := filepath.Join(
		workplaceDir,
		"config.local.json",
	)

	if err := os.WriteFile(
		localSource,
		[]byte(localJSON),
		0644,
	); err != nil {
		t.Fatal(err)
	}

	if _, err := MigrateConfig(dir, false); err != nil {
		t.Fatal(err)
	}

	actual, err := os.ReadFile(localSource)
	if err != nil {
		t.Fatal(err)
	}

	if string(actual) != localJSON {
		t.Fatalf("legacy local JSON was modified")
	}
}

func assertYAMLEqual(
	t *testing.T,
	expected interface{},
	actual interface{},
) {
	t.Helper()

	expectedData, err := yaml.Marshal(expected)
	if err != nil {
		t.Fatal(err)
	}

	actualData, err := yaml.Marshal(actual)
	if err != nil {
		t.Fatal(err)
	}

	if string(expectedData) != string(actualData) {
		t.Errorf(
			"YAML mismatch:\nexpected:\n%s\nactual:\n%s",
			expectedData,
			actualData,
		)
	}
}

func assertNoTemporaryFiles(t *testing.T, dir string) {
	t.Helper()

	patterns := []string{
		filepath.Join(dir, ".workplace-*.yaml.tmp"),
		filepath.Join(dir, ".workplace", "config.local-*.yaml.tmp"),
	}

	for _, pattern := range patterns {
		files, err := filepath.Glob(pattern)
		if err != nil {
			t.Fatal(err)
		}

		if len(files) > 0 {
			t.Errorf("temporary files were not cleaned up: %v", files)
		}
	}
}
