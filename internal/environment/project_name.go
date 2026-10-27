package environment

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

func ResolveProjectName(projectDir string) (string, error) {
	if strings.TrimSpace(projectDir) == "" {
		projectDir = "."
	}

	root, err := filepath.Abs(projectDir)
	if err != nil {
		return "", fmt.Errorf("resolve project directory: %w", err)
	}

	info, err := os.Stat(root)
	if err != nil {
		return "", fmt.Errorf("inspect project directory: %w", err)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("project path is not a directory: %s", root)
	}

	// First: .workplace.yml
	yamlPath := filepath.Join(root, ".workplace.yml")
	data, err := os.ReadFile(yamlPath)
	if err == nil {
		var config struct {
			Name string `yaml:"name"`
		}
		if err := yaml.Unmarshal(data, &config); err != nil {
			return "", fmt.Errorf("parse %s: %w", yamlPath, err)
		}
		if name := strings.TrimSpace(config.Name); name != "" {
			return name, nil
		}
	} else if !os.IsNotExist(err) {
		return "", fmt.Errorf("read %s: %w", yamlPath, err)
	}

	// Second: package.json
	packagePath := filepath.Join(root, "package.json")
	data, err = os.ReadFile(packagePath)
	if err == nil {
		var manifest struct {
			Name string `json:"name"`
		}
		if err := json.Unmarshal(data, &manifest); err != nil {
			return "", fmt.Errorf("parse %s: %w", packagePath, err)
		}
		if name := strings.TrimSpace(manifest.Name); name != "" {
			return name, nil
		}
	} else if !os.IsNotExist(err) {
		return "", fmt.Errorf("read %s: %w", packagePath, err)
	}

	// Third: current project directory name
	name := filepath.Base(filepath.Clean(root))
	if name == "" || name == "." || name == string(filepath.Separator) {
		return "", fmt.Errorf("could not determine project name from %s", root)
	}

	return name, nil
}
