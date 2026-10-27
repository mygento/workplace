package config

import (
	"fmt"
	"path/filepath"
	"strings"
)

type Config struct {
	ProjectDir string             `yaml:"-"`
	Name       string             `yaml:"name,omitempty"`
	Type       string             `yaml:"type"`
	Services   map[string]Service `yaml:"services"`
}

type Service struct {
	Image   string            `yaml:"image,omitempty"`
	Port    *int              `yaml:"port,omitempty"`
	Env     map[string]string `yaml:"env,omitempty"`
	Managed *bool             `yaml:"managed,omitempty"`
}

// IsManaged reports whether Workplace should manage this service locally.
// Services are managed by default.
func (s Service) IsManaged() bool {
	return s.Managed == nil || *s.Managed
}

func (c *Config) Service(name string) (Service, bool) {
	service, ok := c.Services[name]
	return service, ok
}

func (c *Config) HasService(name string) bool {
	_, ok := c.Services[name]
	return ok
}

func normalizeType(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func intPtr(value int) *int {
	return &value
}

func defaults(projectType string) (Config, error) {
	var services map[string]Service

	switch normalizeType(projectType) {
	case "magento2":
		services = map[string]Service{
			"php": {
				Image: "php:8.3-fpm",
			},
			"nginx": {
				Image: "nginx:alpine",
				Port:  intPtr(8080),
			},
			"mysql": {
				Image: "mysql:8.4",
				Env: map[string]string{
					"MYSQL_DATABASE":      "magento",
					"MYSQL_USER":          "magento",
					"MYSQL_PASSWORD":      "magento",
					"MYSQL_ROOT_PASSWORD": "root",
				},
			},
			"opensearch": {
				Image: "opensearchproject/opensearch:latest",
				Env: map[string]string{
					"discovery.type":          "single-node",
					"DISABLE_SECURITY_PLUGIN": "true",
				},
			},
		}

	case "symfony":
		services = map[string]Service{
			"php": {
				Image: "php:8.3-fpm",
			},
			"nginx": {
				Image: "nginx:alpine",
				Port:  intPtr(8080),
			},
		}

	default:
		return Config{}, fmt.Errorf(
			"unknown project type %q",
			projectType,
		)
	}

	return Config{
		Type:     normalizeType(projectType),
		Services: services,
	}, nil
}

func configPath(projectDir string) string {
	return filepath.Join(projectDir, ".workplace.yaml")
}

func localConfigPath(projectDir string) string {
	return filepath.Join(
		projectDir,
		".workplace",
		"config.local.yaml",
	)
}
