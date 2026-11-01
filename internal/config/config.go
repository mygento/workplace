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
	Image         string            `yaml:"image,omitempty"`
	Port          *int              `yaml:"port,omitempty"`
	Env           map[string]string `yaml:"env,omitempty"`
	Volumes       []Volume          `yaml:"volumes,omitempty"`
	DependsOn     []string          `yaml:"depends_on,omitempty"`
	ContainerName string            `yaml:"container_name,omitempty"`
	Managed       *bool             `yaml:"managed,omitempty"`
}

type Volume struct {
	Name     string `yaml:"name,omitempty"`
	Source   string `yaml:"source,omitempty"`
	Target   string `yaml:"target"`
	Type     string `yaml:"type,omitempty"`
	ReadOnly bool   `yaml:"read_only,omitempty"`
}

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

var optionalServices = map[string]struct{}{
	"redis":    {},
	"elastic":  {},
	"varnish":  {},
	"rabbitmq": {},
}

func IsOptionalService(name string) bool {
	_, ok := optionalServices[name]
	return ok
}

func serviceDefaults(projectType string, projectName string, name string) (Service, bool) {
	projectMount := "/var/www/" + projectType

	switch name {
	case "php":
		return Service{
			Image: "php:8.3-fpm",
			Env: map[string]string{
				"PHPFPM_USER": "${USERID}",
			},
			DependsOn: []string{"db"},
			Volumes: []Volume{
				{
					Source: ".",
					Target: projectMount,
					Type:   "bind",
				},
			},
		}, true

	case "nginx":
		return Service{
			Image: "nginx:alpine",
			Port:  intPtr(8081),
			Env: map[string]string{
				"NGINX_USER": "${USERID}",
			},
			DependsOn: []string{"php"},
			Volumes: []Volume{
				{
					Source:   ".",
					Target:   projectMount,
					Type:     "bind",
					ReadOnly: true,
				},
				{
					Source:   ".workplace/runtime/nginx",
					Target:   "/etc/nginx/sites-enabled",
					Type:     "bind",
					ReadOnly: true,
				},
			},
		}, true

	case "mysql":
		return Service{
			Image: "mysql:8.4",
			Port:  intPtr(3306),
			Env: map[string]string{
				"MYSQL_DATABASE":      projectType,
				"MYSQL_PASSWORD":      "mygento",
				"MYSQL_ROOT_PASSWORD": "mygento",
				"MYSQL_USER":          "mygento",
			},
			Volumes: []Volume{
				{
					Name:   projectName + "-db",
					Target: "/var/lib/mysql",
					Type:   "volume",
				},
			},
		}, true

	case "redis":
		return Service{
			Image: "redis:alpine",
		}, true

	case "elastic":
		return Service{
			Image: "opensearchproject/opensearch:3",
			Env: map[string]string{
				"discovery.type":          "single-node",
				"DISABLE_SECURITY_PLUGIN": "true",
			},
			Volumes: []Volume{
				{
					Name:   projectName + "-elastic",
					Target: "/usr/share/opensearch/data",
					Type:   "volume",
				},
			},
		}, true

	case "varnish":
		return Service{
			Image: "varnish:stable",
			Port:  intPtr(8081),
		}, true

	case "rabbitmq":
		return Service{
			Image: "rabbitmq:management",
			Port:  intPtr(5672),
			Volumes: []Volume{
				{
					Name:   "rabbitmq",
					Target: "/var/lib/rabbitmq",
					Type:   "volume",
				},
			},
		}, true

	default:
		return Service{}, false
	}
}

func loadDefaults(projectType string, projectName string) (Config, error) {
	projectType = normalizeType(projectType)
	services := make(map[string]Service)
	loadCommonServices(projectType, projectName, services)

	switch projectType {
	case "magento2":
		service, _ := serviceDefaults(projectType, projectName, "elastic")
		services["elastic"] = service
		return Config{
			Type:     projectType,
			Services: services,
		}, nil
	case "symfony":
		return Config{
			Type:     projectType,
			Services: services,
		}, nil

	default:
		return Config{}, fmt.Errorf(
			"unknown project type %q",
			projectType,
		)
	}
}

func loadCommonServices(projectType string, projectName string, services map[string]Service) {
	for _, name := range []string{"php", "nginx", "mysql"} {
		service, _ := serviceDefaults(projectType, projectName, name)
		if name == "mysql" {
			service.ContainerName = projectName + "-db"
		} else {
			service.ContainerName = projectName + "-" + name
		}
		services[name] = service
	}
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
