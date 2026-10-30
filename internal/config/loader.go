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

	projectPath := configPath(root)

	projectConfig, err := readConfig(projectPath, false)
	if err != nil {
		return Config{}, err
	}
	projectConfig.ProjectDir = root

	projectType := normalizeType(projectConfig.Type)
	if projectType == "" {
		return Config{}, fmt.Errorf(
			"project type is missing in %s",
			projectPath,
		)
	}

	if err := validateManagedOverrides(projectConfig, projectPath); err != nil {
		return Config{}, err
	}

	cfg, err := loadDefaults(projectType, projectConfig.Name)
	if err != nil {
		return Config{}, err
	}

	if err := mergeConfig(&cfg, projectConfig, false); err != nil {
		return Config{}, fmt.Errorf("merge project config: %w", err)
	}

	localPath := localConfigPath(root)

	localConfig, err := readConfig(localPath, true)
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

	if err := mergeConfig(&cfg, localConfig, true); err != nil {
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

// validateManagedOverrides prevents shared project configuration from
// controlling local native service lifecycle.
func validateManagedOverrides(cfg Config, path string) error {
	for name, service := range cfg.Services {
		if service.Managed != nil {
			return fmt.Errorf(
				"services.%s.managed is only allowed in %s",
				name,
				localConfigPath(filepath.Dir(path)),
			)
		}
	}

	return nil
}

// mergeConfig merges project or local overrides into the resolved
// configuration. Optional services are included only if declared
// in the source configuration.
func mergeConfig(dst *Config, src Config, allowManaged bool) error {
	for name, override := range src.Services {
		service, exists := dst.Services[name]

		if !exists {
			if !IsOptionalService(name) {
				return fmt.Errorf(
					"unknown service %q in configuration",
					name,
				)
			}

			var ok bool
			service, ok = serviceDefaults(dst.Type, dst.Name, name)
			if !ok {
				return fmt.Errorf(
					"no defaults defined for service %q",
					name,
				)
			}
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

		if override.Volumes != nil {
			service.Volumes = override.Volumes
		}

		if override.DependsOn != nil {
			service.DependsOn = override.DependsOn
		}

		if override.Managed != nil {
			if !allowManaged {
				return fmt.Errorf(
					"services.%s.managed is only allowed in local configuration",
					name,
				)
			}

			managed := *override.Managed
			service.Managed = &managed
		}

		dst.Services[name] = service
	}

	if src.Name != "" {
		dst.Name = src.Name
	}
	if src.ProjectDir != "" {
		dst.ProjectDir = src.ProjectDir
	}

	return nil
}
