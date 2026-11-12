package environment

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

var (
	ErrLegacyConfigNotFound = errors.New("workplace configuration not found in package.json")
	ErrTargetConfigExists   = errors.New(".workplace.yaml already exists")
	ErrLocalConfigExists    = errors.New(".workplace/config.local.yaml already exists")
)

type MigrationResult struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Type        string `json:"type"`
	Written     bool   `json:"written"`
}

type legacyManifest struct {
	Name      string          `json:"name"`
	Workplace json.RawMessage `json:"workplace"`
}

func MigrateConfig(projectDir string, overwrite bool) (result *MigrationResult, retErr error) {
	if projectDir == "" {
		projectDir = "."
	}

	root, err := filepath.Abs(projectDir)
	if err != nil {
		return nil, fmt.Errorf("resolve project directory: %w", err)
	}

	info, err := os.Stat(root)
	if err != nil {
		return nil, fmt.Errorf("project directory: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("project path is not a directory: %s", root)
	}

	source := filepath.Join(root, "package.json")
	destination := filepath.Join(root, ".workplace.yaml")

	data, err := os.ReadFile(source)
	if err != nil {
		return nil, fmt.Errorf("read package.json: %w", err)
	}

	var manifest legacyManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, fmt.Errorf("parse package.json: %w", err)
	}

	if len(manifest.Workplace) == 0 ||
		strings.TrimSpace(string(manifest.Workplace)) == "null" {
		return nil, ErrLegacyConfigNotFound
	}

	var workplace map[string]json.RawMessage
	if err := json.Unmarshal(manifest.Workplace, &workplace); err != nil {
		return nil, fmt.Errorf("parse package.json workplace config: %w", err)
	}

	var projectType string
	if err := json.Unmarshal(workplace["type"], &projectType); err != nil {
		return nil, errors.New("legacy workplace configuration has no type")
	}

	projectType = strings.TrimSpace(projectType)
	if projectType == "" {
		return nil, errors.New("legacy workplace configuration has no type")
	}

	projectName := strings.TrimSpace(manifest.Name)
	if projectName == "" {
		projectName = filepath.Base(root)
	}

	// Build the new project configuration.
	services, err := migrateServices(workplace, "type", projectType)
	if err != nil {
		return nil, fmt.Errorf("migrate project services: %w", err)
	}

	config := map[string]interface{}{
		"name":     projectName,
		"type":     projectType,
		"services": services,
	}

	yamlData, err := yaml.Marshal(config)
	if err != nil {
		return nil, fmt.Errorf("encode YAML: %w", err)
	}

	// Convert local configuration when the legacy file exists.
	localSource := filepath.Join(root, ".workplace", "config.local.json")
	localDestination := filepath.Join(root, ".workplace", "config.local.yaml")

	var localYAML []byte

	localData, err := os.ReadFile(localSource)
	switch {
	case err == nil:
		var localConfig map[string]json.RawMessage

		if err := json.Unmarshal(localData, &localConfig); err != nil {
			return nil, fmt.Errorf("parse local config: %w", err)
		}

		localServices, err := migrateServices(localConfig)
		if err != nil {
			return nil, fmt.Errorf("migrate local services: %w", err)
		}

		localYAML, err = yaml.Marshal(map[string]interface{}{
			"services": localServices,
		})
		if err != nil {
			return nil, fmt.Errorf("encode local YAML: %w", err)
		}

	case errors.Is(err, os.ErrNotExist):
		// Local configuration is optional.

	default:
		return nil, fmt.Errorf("read local config: %w", err)
	}

	// Check both destinations before writing either file.
	if err := checkMigrationDestination(destination, overwrite, ErrTargetConfigExists); err != nil {
		return nil, err
	}

	if localYAML != nil {
		if err := checkMigrationDestination(localDestination, overwrite, ErrLocalConfigExists); err != nil {
			return nil, err
		}
	}

	// Write the project configuration atomically.
	if err := writeMigrationFile(root, ".workplace-*.yaml.tmp", destination, yamlData, overwrite, ErrTargetConfigExists); err != nil {
		return nil, fmt.Errorf("write project config: %w", err)
	}

	// Write the local configuration atomically using the same implementation.
	if localYAML != nil {
		if err := writeMigrationFile(
			filepath.Dir(localDestination),
			"config.local-*.yaml.tmp",
			localDestination,
			localYAML,
			overwrite,
			ErrLocalConfigExists,
		); err != nil {
			return nil, fmt.Errorf("write local config: %w", err)
		}
	}

	return &MigrationResult{
		Source:      source,
		Destination: destination,
		Type:        projectType,
		Written:     true,
	}, nil
}

// migrateServices converts legacy service definitions into the new
// configuration format. It is shared by project and local migrations.
func migrateServices(
	legacy map[string]json.RawMessage,
	excluded ...string,
) (map[string]interface{}, error) {
	skip := make(map[string]struct{}, len(excluded))

	for _, name := range excluded {
		skip[name] = struct{}{}
	}

	services := make(map[string]interface{})

	for name, raw := range legacy {
		if _, ok := skip[name]; ok {
			continue
		}

		if len(raw) == 0 || strings.TrimSpace(string(raw)) == "null" {
			continue
		}

		var service interface{}
		if err := json.Unmarshal(raw, &service); err != nil {
			return nil, fmt.Errorf(
				"parse legacy %s service configuration: %w",
				name,
				err,
			)
		}

		// The new configuration uses "elastic" instead of
		// the legacy "elasticsearch" service name.
		serviceName := name
		if name == "elasticsearch" {
			serviceName = "elastic"
		}

		services[serviceName] = service
	}

	return services, nil
}

func checkMigrationDestination(
	destination string,
	overwrite bool,
	existsError error,
) error {
	if _, err := os.Stat(destination); err == nil && !overwrite {
		return existsError
	} else if err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("check destination: %w", err)
	}

	return nil
}

// writeMigrationFile uses the same atomic write strategy for both
// project and local configuration files.
func writeMigrationFile(
	dir string,
	pattern string,
	destination string,
	data []byte,
	overwrite bool,
	existsError error,
) (retErr error) {
	tmp, err := os.CreateTemp(dir, pattern)
	if err != nil {
		return fmt.Errorf("create temporary config: %w", err)
	}

	tmpName := tmp.Name()
	closed := false

	defer func() {
		if !closed {
			if err := tmp.Close(); err != nil {
				retErr = errors.Join(
					retErr,
					fmt.Errorf("close temporary config: %w", err),
				)
			}
		}

		if err := os.Remove(tmpName); err != nil && !errors.Is(err, os.ErrNotExist) {
			retErr = errors.Join(
				retErr,
				fmt.Errorf("remove temporary config: %w", err),
			)
		}
	}()

	if _, err := tmp.Write(data); err != nil {
		return fmt.Errorf("write temporary config: %w", err)
	}

	if err := tmp.Sync(); err != nil {
		return fmt.Errorf("sync temporary config: %w", err)
	}

	if err := tmp.Close(); err != nil {
		closed = true
		return fmt.Errorf("close temporary config: %w", err)
	}
	closed = true

	if !overwrite {
		if err := os.Link(tmpName, destination); err != nil {
			if errors.Is(err, os.ErrExist) {
				return existsError
			}
			return fmt.Errorf("create destination config: %w", err)
		}
	} else {
		if err := os.Rename(tmpName, destination); err != nil {
			return fmt.Errorf("replace destination config: %w", err)
		}
	}

	return nil
}
