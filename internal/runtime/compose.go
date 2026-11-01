package runtime

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/mygento/workplace/internal/config"
	"gopkg.in/yaml.v3"
)

type composeDocument struct {
	Name     string                    `yaml:"name"`
	Services map[string]composeService `yaml:"services"`
	Volumes  map[string]composeVolume  `yaml:"volumes,omitempty"`
}

type composeService struct {
	ContainerName string         `yaml:"container_name,omitempty"`
	Image         string         `yaml:"image,omitempty"`
	Environment   []quotedString `yaml:"environment,omitempty"`
	Ports         []quotedString `yaml:"ports,omitempty"`
	Volumes       []string       `yaml:"volumes,omitempty"`
	DependsOn     []string       `yaml:"depends_on,omitempty"`
}

type quotedString string

func (s quotedString) MarshalYAML() (interface{}, error) {
	return &yaml.Node{
		Kind:  yaml.ScalarNode,
		Tag:   "!!str",
		Value: string(s),
		Style: yaml.DoubleQuotedStyle,
	}, nil
}

type composeVolume struct{}

func composeServiceName(name string) string {
	switch name {
	case "mysql":
		return "db"
	case "elasticsearch", "opensearch":
		return "elastic"
	default:
		return name
	}
}

func composeContainerPort(name string) (int, bool) {
	switch name {
	case "nginx":
		return 80, true
	case "mysql":
		return 3306, true
	case "redis":
		return 6379, true
	case "elasticsearch":
		return 9200, true
	case "varnish":
		return 80, true
	case "rabbitmq":
		return 5672, true
	default:
		return 0, false
	}
}

func composeVolumeMount(
	projectDir string,
	volume config.Volume,
) (string, error) {
	target := strings.TrimSpace(volume.Target)
	if target == "" {
		return "", fmt.Errorf("volume target cannot be empty")
	}

	var source string

	switch volume.Type {
	case "volume":
		source = strings.TrimSpace(volume.Name)
		if source == "" {
			return "", fmt.Errorf(
				"named volume for target %q has no name",
				target,
			)
		}

	case "bind", "":
		source = strings.TrimSpace(volume.Source)
		if source == "" {
			source = projectDir
		} else if !filepath.IsAbs(source) {
			source = filepath.Join(projectDir, source)
		}

		absolute, err := filepath.Abs(source)
		if err != nil {
			return "", fmt.Errorf("resolve volume source: %w", err)
		}
		source = absolute

	default:
		return "", fmt.Errorf(
			"unsupported volume type %q",
			volume.Type,
		)
	}

	mount := source + ":" + target
	if volume.ReadOnly {
		mount += ":ro"
	}

	return mount, nil
}

func renderCompose(cfg config.Config) ([]byte, error) {
	projectDir, err := filepath.Abs(cfg.ProjectDir)
	if err != nil {
		return nil, fmt.Errorf("resolve project directory: %w", err)
	}

	document := composeDocument{
		Name:     dockerProjectName(cfg),
		Services: make(map[string]composeService),
		Volumes:  make(map[string]composeVolume),
	}

	serviceNames := make([]string, 0, len(cfg.Services))
	for name := range cfg.Services {
		serviceNames = append(serviceNames, name)
	}

	// Stable ordering makes generated YAML easier to review.
	sortStrings(serviceNames)

	for _, name := range serviceNames {
		service := cfg.Services[name]

		// Optional services without an image are not enabled.
		if strings.TrimSpace(service.Image) == "" {
			if config.IsOptionalService(name) {
				continue
			}

			return nil, fmt.Errorf(
				"service %q has no image",
				name,
			)
		}

		composeName := composeServiceName(name)

		entry := composeService{
			Image: service.Image,
		}
		if service.ContainerName != "" {
			entry.ContainerName = service.ContainerName
		}

		if len(service.Env) > 0 {
			keys := make([]string, 0, len(service.Env))
			for key := range service.Env {
				keys = append(keys, key)
			}
			sort.Strings(keys)

			entry.Environment = make([]quotedString, 0, len(service.Env))

			for _, key := range keys {
				entry.Environment = append(
					entry.Environment,
					quotedString(key+"="+service.Env[key]),
				)
			}
		}

		if service.Port != nil {
			containerPort, ok := composeContainerPort(name)
			if !ok {
				return nil, fmt.Errorf(
					"cannot determine container port for service %q",
					name,
				)
			}

			entry.Ports = []quotedString{
				quotedString(
					strconv.Itoa(*service.Port) + ":" +
						strconv.Itoa(containerPort),
				),
			}
		}

		for _, dependency := range service.DependsOn {
			dependency = composeServiceName(dependency)

			// Dependencies are only emitted when the target service
			// is present in the generated Compose document.
			if _, exists := cfg.Services[dependency]; !exists {
				originalName := dependency
				for configName := range cfg.Services {
					if composeServiceName(configName) == dependency {
						originalName = configName
						break
					}
				}

				if _, exists := cfg.Services[originalName]; !exists {
					continue
				}
			}

			entry.DependsOn = append(entry.DependsOn, dependency)
		}

		for _, volume := range service.Volumes {
			mount, err := composeVolumeMount(projectDir, volume)
			if err != nil {
				return nil, fmt.Errorf(
					"service %q: %w",
					name,
					err,
				)
			}

			entry.Volumes = append(entry.Volumes, mount)

			if volume.Type == "volume" {
				document.Volumes[volume.Name] = composeVolume{}
			}
		}

		document.Services[composeName] = entry
	}

	if len(document.Services) == 0 {
		return nil, fmt.Errorf("no services configured for Docker Compose")
	}

	if len(document.Volumes) == 0 {
		document.Volumes = nil
	}

	data, err := yaml.Marshal(document)
	if err != nil {
		return nil, fmt.Errorf("marshal Compose configuration: %w", err)
	}

	return data, nil
}

func writeCompose(cfg config.Config, path string) error {
	data, err := renderCompose(cfg)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("create Compose directory: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("write Compose configuration: %w", err)
	}

	return nil
}

func sortStrings(values []string) {
	for i := 1; i < len(values); i++ {
		for j := i; j > 0 && values[j] < values[j-1]; j-- {
			values[j], values[j-1] = values[j-1], values[j]
		}
	}
}
