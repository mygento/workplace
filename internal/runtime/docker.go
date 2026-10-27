package runtime

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"

	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/flags"
	"github.com/docker/compose/v5/cmd/display"
	"github.com/docker/compose/v5/pkg/api"
	"github.com/docker/compose/v5/pkg/compose"
	"github.com/mygento/workplace/internal/config"
	"gopkg.in/yaml.v3"
)

type DockerRuntime struct {
	mu      sync.Mutex
	service api.Compose
	project string
	loaded  bool
}

type composeConfig struct {
	Name     string                          `yaml:"name"`
	Services map[string]composeConfigService `yaml:"services"`
}

type composeConfigService struct {
	Image       string            `yaml:"image"`
	Environment map[string]string `yaml:"environment,omitempty"`
	Ports       []string          `yaml:"ports,omitempty"`
}

func NewDocker() *DockerRuntime {
	return &DockerRuntime{}
}

var invalidProjectName = regexp.MustCompile(`[^a-z0-9_-]+`)

func dockerProjectName(cfg config.Config) string {
	name := strings.ToLower(strings.TrimSpace(cfg.Name))
	if name == "" {
		name = filepath.Base(cfg.ProjectDir)
	}

	name = invalidProjectName.ReplaceAllString(name, "-")
	name = strings.Trim(name, "-_")
	if name == "" {
		name = "workplace"
	}
	if name[0] < 'a' || name[0] > 'z' {
		name = "workplace-" + name
	}

	return name
}

func (r *DockerRuntime) Start(
	ctx context.Context,
	cfg config.Config,
) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.loaded {
		return nil
	}

	fmt.Println("Initializing Docker runtime...")

	dockerCLI, err := command.NewDockerCli()
	if err != nil {
		return fmt.Errorf("create Docker client: %w", err)
	}

	if err := dockerCLI.Initialize(&flags.ClientOptions{}); err != nil {
		return fmt.Errorf("initialize Docker client: %w", err)
	}

	fmt.Println("Initializing Docker Compose SDK...")

	service, err := compose.NewComposeService(
		dockerCLI,
		compose.WithEventProcessor(display.Plain(os.Stdout)),
	)
	if err != nil {
		return fmt.Errorf("initialize compose sdk: %w", err)
	}

	projectName := dockerProjectName(cfg)

	services := make(map[string]composeConfigService)
	names := make([]string, 0, len(cfg.Services))

	for name := range cfg.Services {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		svc := cfg.Services[name]

		if !svc.IsManaged() {
			fmt.Printf("Skipping unmanaged service %q\n", name)
			continue
		}

		if strings.TrimSpace(svc.Image) == "" {
			return fmt.Errorf("service %q has no image", name)
		}

		entry := composeConfigService{
			Image:       svc.Image,
			Environment: svc.Env,
		}

		if svc.Port != nil {
			containerPort, ok := dockerContainerPort(name)
			if !ok {
				return fmt.Errorf(
					"cannot infer container port for service %q; add explicit container-port configuration",
					name,
				)
			}

			entry.Ports = []string{
				fmt.Sprintf("%d:%d", *svc.Port, containerPort),
			}
		}

		services[name] = entry
	}

	if len(services) == 0 {
		return fmt.Errorf("no managed services to start")
	}

	runtimeDir := filepath.Join(
		cfg.ProjectDir,
		".workplace",
		"runtime",
	)

	if err := os.MkdirAll(runtimeDir, 0755); err != nil {
		return fmt.Errorf("create runtime directory: %w", err)
	}

	composePath := filepath.Join(runtimeDir, "compose.yaml")

	data, err := yaml.Marshal(composeConfig{
		Name:     projectName,
		Services: services,
	})
	if err != nil {
		return fmt.Errorf("marshal Compose configuration: %w", err)
	}

	if err := os.WriteFile(composePath, data, 0600); err != nil {
		return fmt.Errorf("write Compose configuration: %w", err)
	}

	fmt.Printf("Loading Docker Compose project %q...\n", projectName)

	project, err := service.LoadProject(ctx, api.ProjectLoadOptions{
		ConfigPaths: []string{composePath},
		ProjectName: projectName,
	})
	if err != nil {
		return fmt.Errorf("load Compose project: %w", err)
	}

	fmt.Println("Starting Docker services...")

	if err := service.Up(ctx, project, api.UpOptions{
		Create: api.CreateOptions{},
		Start:  api.StartOptions{},
	}); err != nil {
		return fmt.Errorf(
			"start Compose project %q: %w",
			projectName,
			err,
		)
	}

	r.service = service
	r.project = projectName
	r.loaded = true

	fmt.Printf("Started Docker Compose project %q\n", projectName)

	return nil
}

func (r *DockerRuntime) Stop() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.loaded {
		return nil
	}

	fmt.Printf("Stopping Docker Compose project %q...\n", r.project)

	if err := r.service.Down(
		context.Background(),
		r.project,
		api.DownOptions{},
	); err != nil {
		return fmt.Errorf(
			"stop Compose project %q: %w",
			r.project,
			err,
		)
	}

	r.loaded = false
	r.project = ""
	r.service = nil

	fmt.Println("Docker Compose project stopped")

	return nil
}

func dockerContainerPort(name string) (int, bool) {
	port, ok := map[string]int{
		"nginx":      80,
		"php":        9000,
		"mysql":      3306,
		"opensearch": 9200,
	}[name]

	return port, ok
}
