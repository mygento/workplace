package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Loader struct {
	ProjectDir string
}

func NewLoader(projectDir string) *Loader {
	return &Loader{
		ProjectDir: projectDir,
	}
}

func (l *Loader) Load() (Config, error) {
	root, err := filepath.Abs(l.ProjectDir)
	if err != nil {
		return Config{}, fmt.Errorf(
			"resolve project directory: %w",
			err,
		)
	}

	projectConfig, err := readConfig(configPath(root), false)
	if err != nil {
		return Config{}, err
	}

	projectType := normalizeType(projectConfig.Type)
	if projectType == "" {
		return Config{}, fmt.Errorf(
			"project type is missing in %s",
			configPath(root),
		)
	}

	cfg, err := defaults(projectType)
	if err != nil {
		return Config{}, err
	}

	if err := mergeConfig(&cfg, projectConfig); err != nil {
		return Config{}, fmt.Errorf("merge project config: %w", err)
	}

	localConfig, err := readConfig(localConfigPath(root), true)
	if err != nil {
		return Config{}, err
	}

	if localConfig.Type != "" &&
		normalizeType(localConfig.Type) != cfg.Type {
		return Config{}, fmt.Errorf(
			"local configuration cannot change project type %q",
			cfg.Type,
		)
	}
	if localConfig.Name != "" {
		return Config{}, fmt.Errorf(
			"local configuration cannot change project name %q",
			cfg.Name,
		)
	}

	if err := mergeConfig(&cfg, localConfig); err != nil {
		return Config{}, fmt.Errorf("merge local config: %w", err)
	}

	cfg.ProjectDir = root
	return cfg, nil
}

func readConfig(path string, optional bool) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if optional && os.IsNotExist(err) {
			return Config{}, nil
		}

		return Config{}, fmt.Errorf("read %s: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("parse %s: %w", path, err)
	}

	return cfg, nil
}

func mergeConfig(dst *Config, src Config) error {
	for name, override := range src.Services {
		service, exists := dst.Services[name]
		if !exists {
			service = Service{}
		}

		if override.Image != "" {
			service.Image = override.Image
		}

		if override.Port != nil {
			port := *override.Port
			service.Port = &port
		}

		if override.Env != nil {
			if service.Env == nil {
				service.Env = make(map[string]string)
			}

			for key, value := range override.Env {
				service.Env[key] = value
			}
		}

		dst.Services[name] = service
	}

	return nil
}
