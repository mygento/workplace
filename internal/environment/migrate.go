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
	ErrTargetConfigExists   = errors.New(".workplace.yml already exists")
)

type MigrationResult struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Type        string `json:"type"`
	Written     bool   `json:"written"`
}

type legacyManifest struct {
	Workplace json.RawMessage `json:"workplace"`
}

type legacyWorkplaceConfig struct {
	Type          string                 `json:"type"`
	PHP           map[string]interface{} `json:"php,omitempty"`
	Nginx         map[string]interface{} `json:"nginx,omitempty"`
	MySQL         map[string]interface{} `json:"mysql,omitempty"`
	Elasticsearch map[string]interface{} `json:"elasticsearch,omitempty"`
}

// MigrateConfig extracts package.json's workplace property and writes it
// as .workplace.yml. Existing destination files are never overwritten
// unless overwrite is explicitly true.
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
	destination := filepath.Join(root, ".workplace.yml")

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

	var config legacyWorkplaceConfig
	if err := json.Unmarshal(manifest.Workplace, &config); err != nil {
		return nil, fmt.Errorf("parse package.json workplace config: %w", err)
	}
	if strings.TrimSpace(config.Type) == "" {
		return nil, errors.New("legacy workplace configuration has no type")
	}

	if _, err := os.Stat(destination); err == nil && !overwrite {
		return nil, ErrTargetConfigExists
	} else if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("check destination: %w", err)
	}

	yamlData, err := yaml.Marshal(config)
	if err != nil {
		return nil, fmt.Errorf("encode YAML: %w", err)
	}

	// Write atomically so an interrupted migration does not leave a
	// partially written configuration file.
	tmp, err := os.CreateTemp(root, ".workplace-*.yml.tmp")
	if err != nil {
		return nil, fmt.Errorf("create temporary config: %w", err)
	}

	tmpName := tmp.Name()
	closed := false

	defer func() {
		if !closed {
			if err := tmp.Close(); err != nil {
				retErr = errors.Join(retErr, fmt.Errorf("close temporary config: %w", err))
			}
		}

		if err := os.Remove(tmpName); err != nil && !os.IsNotExist(err) {
			retErr = errors.Join(retErr, fmt.Errorf("remove temporary config: %w", err))
		}
	}()

	if _, err := tmp.Write(yamlData); err != nil {
		return nil, fmt.Errorf("write temporary config: %w", err)
	}

	if err := tmp.Sync(); err != nil {
		return nil, fmt.Errorf("sync temporary config: %w", err)
	}

	if err := tmp.Close(); err != nil {
		closed = true
		return nil, fmt.Errorf("close temporary config: %w", err)
	}
	closed = true

	if !overwrite {
		// A hard link provides create-if-absent behavior, preventing
		// accidental replacement if another process creates the file.
		if err := os.Link(tmpName, destination); err != nil {
			if os.IsExist(err) {
				return nil, ErrTargetConfigExists
			}
			return nil, fmt.Errorf("create destination config: %w", err)
		}
	} else {
		if err := os.Rename(tmpName, destination); err != nil {
			return nil, fmt.Errorf("replace destination config: %w", err)
		}
	}

	return &MigrationResult{
		Source:      source,
		Destination: destination,
		Type:        config.Type,
		Written:     true,
	}, nil
}
