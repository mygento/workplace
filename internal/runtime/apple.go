package runtime

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/mygento/workplace/internal/config"
	"github.com/mygento/workplace/internal/project"
)

type AppleRuntime struct{}

type appleNetwork struct {
	ID            string `json:"id"`
	Configuration struct {
		Name string `json:"name"`
	} `json:"configuration"`
}

func NewApple() *AppleRuntime {
	return &AppleRuntime{}
}

func appleCommand(
	ctx context.Context,
	args ...string,
) ([]byte, error) {
	fmt.Printf(
		"[apple] $ container %s\n",
		shellCommand("container", args),
	)

	return RunCommand(ctx, "container", args...)
}

func appleCommandJSON[T any](
	ctx context.Context,
	args ...string,
) (T, error) {
	var result T

	jsonArgs := append(
		append([]string(nil), args...),
		"--format",
		"json",
	)

	output, err := appleCommand(ctx, jsonArgs...)
	if err != nil {
		return result, err
	}

	if err := json.Unmarshal(output, &result); err != nil {
		return result, fmt.Errorf(
			"decode Apple Container JSON: %w; output: %s",
			err,
			strings.TrimSpace(string(output)),
		)
	}

	return result, nil
}

func (r *AppleRuntime) Start(
	ctx context.Context,
	cfg config.Config,
) error {
	if cfg.ProjectDir == "" {
		return fmt.Errorf("project directory is not configured")
	}

	if _, err := exec.LookPath("container"); err != nil {
		return fmt.Errorf(
			"apple container CLI not found: %w",
			err,
		)
	}

	if _, err := os.Stat(cfg.ProjectDir); err != nil {
		return fmt.Errorf(
			"access project directory %q: %w",
			cfg.ProjectDir,
			err,
		)
	}

	output, err := appleCommand(ctx, "system", "status")
	if err != nil {
		return fmt.Errorf(
			"apple container system is not running; run `container system start`: %w: %s",
			err,
			strings.TrimSpace(string(output)),
		)
	}

	names, err := appleServiceOrder(cfg)
	if err != nil {
		return fmt.Errorf(
			"resolve service dependencies: %w",
			err,
		)
	}

	if len(names) == 0 {
		return fmt.Errorf("no managed services configured")
	}

	for _, name := range names {
		if strings.TrimSpace(cfg.Services[name].Image) == "" {
			return fmt.Errorf(
				"service %q has no image",
				name,
			)
		}
	}

	if cfg.HasService("nginx") {
		if err := project.Install(cfg.ProjectDir, cfg.Type); err != nil {
			return fmt.Errorf(
				"install project nginx configuration: %w",
				err,
			)
		}
	}

	projectName := cfg.Name
	network := appleNetworkName(projectName)

	exists, err := appleNetworkExists(ctx, network)
	if err != nil {
		return fmt.Errorf(
			"inspect Apple Container networks: %w",
			err,
		)
	}

	if exists {
		fmt.Printf(
			"[apple] reusing network %q\n",
			network,
		)
	} else {
		output, err := appleCommand(
			ctx,
			"network",
			"create",
			network,
		)

		if err != nil {
			// Handle concurrent network creation.
			exists, checkErr := appleNetworkExists(ctx, network)

			if checkErr != nil {
				return fmt.Errorf(
					"create apple container network %q: %w; verify network: %v",
					network,
					err,
					checkErr,
				)
			}

			if !exists {
				return fmt.Errorf(
					"create apple container network %q: %w: %s",
					network,
					err,
					strings.TrimSpace(string(output)),
				)
			}

			fmt.Printf(
				"[apple] network %q created concurrently; reusing\n",
				network,
			)
		}
	}

	var started []string

	rollback := func(cause error) error {
		var failures []string

		cleanupCtx := context.WithoutCancel(ctx)

		for i := len(started) - 1; i >= 0; i-- {
			id := started[i]

			output, err := appleCommand(
				cleanupCtx,
				"delete",
				"--force",
				id,
			)

			if err != nil && appleContainerExists(cleanupCtx, id) {
				failures = append(
					failures,
					fmt.Sprintf(
						"%s: %v: %s",
						id,
						err,
						strings.TrimSpace(string(output)),
					),
				)
			}
		}

		if len(failures) > 0 {
			return fmt.Errorf(
				"%w; rollback errors: %s",
				cause,
				strings.Join(failures, "; "),
			)
		}

		return cause
	}

	for _, name := range names {
		service := cfg.Services[name]

		id := appleServiceContainerID(
			projectName,
			name,
			service,
		)

		if appleContainerExists(ctx, id) {
			if appleContainerStopped(ctx, id) {
				output, err := appleCommand(
					ctx,
					"start",
					id,
				)

				if err != nil {
					return rollback(fmt.Errorf(
						"start existing apple container service %q: %w: %s",
						name,
						err,
						strings.TrimSpace(string(output)),
					))
				}

				started = append(started, id)
			} else {
				fmt.Printf(
					"[apple] container service %q already running; skipping\n",
					name,
				)
			}

			continue
		}

		args, err := applePrepareService(
			ctx,
			cfg,
			name,
			service,
			network,
		)

		if err != nil {
			return rollback(err)
		}

		output, err := appleCommand(ctx, args...)
		if err != nil {
			return rollback(fmt.Errorf(
				"start apple container service %q: %w: %s",
				name,
				err,
				strings.TrimSpace(string(output)),
			))
		}

		started = append(started, id)

		fmt.Printf(
			"[apple] started container service %q\n",
			name,
		)
	}

	fmt.Printf(
		"[apple] started project %q\n",
		projectName,
	)

	return nil
}

func (r *AppleRuntime) Stop(
	ctx context.Context,
	cfg config.Config,
) error {
	if _, err := exec.LookPath("container"); err != nil {
		return fmt.Errorf(
			"apple container CLI not found: %w",
			err,
		)
	}

	names, err := appleServiceOrder(cfg)
	if err != nil {
		return fmt.Errorf(
			"resolve service dependencies: %w",
			err,
		)
	}

	projectName := cfg.Name

	var failures []string

	// Stop dependants before their dependencies.
	for i := len(names) - 1; i >= 0; i-- {
		name := names[i]
		service := cfg.Services[name]

		id := appleServiceContainerID(
			projectName,
			name,
			service,
		)

		if !appleContainerExists(ctx, id) {
			continue
		}

		if appleContainerStopped(ctx, id) {
			continue
		}

		stopCtx, cancel := context.WithTimeout(
			ctx,
			20*time.Second,
		)

		output, err := appleCommand(
			stopCtx,
			"stop",
			"--time",
			"10",
			id,
		)

		cancel()

		if err != nil && !appleContainerStopped(ctx, id) {
			failures = append(
				failures,
				fmt.Sprintf(
					"%s: %v: %s",
					id,
					err,
					strings.TrimSpace(string(output)),
				),
			)
		}
	}

	if len(failures) > 0 {
		return fmt.Errorf(
			"stop apple container services: %s",
			strings.Join(failures, "; "),
		)
	}

	fmt.Printf(
		"[apple] stopped project %q\n",
		projectName,
	)

	return nil
}

func (r *AppleRuntime) Delete(
	ctx context.Context,
	cfg config.Config,
) error {
	if _, err := exec.LookPath("container"); err != nil {
		return fmt.Errorf(
			"apple container CLI not found: %w",
			err,
		)
	}

	names, err := appleServiceOrder(cfg)
	if err != nil {
		return fmt.Errorf(
			"resolve service dependencies: %w",
			err,
		)
	}

	projectName := cfg.Name

	var failures []string

	// Delete dependants before their dependencies.
	for i := len(names) - 1; i >= 0; i-- {
		name := names[i]
		service := cfg.Services[name]

		id := appleServiceContainerID(
			projectName,
			name,
			service,
		)

		if !appleContainerExists(ctx, id) {
			continue
		}

		output, err := appleCommand(
			ctx,
			"delete",
			"--force",
			id,
		)

		if err != nil && appleContainerExists(ctx, id) {
			failures = append(
				failures,
				fmt.Sprintf(
					"%s: %v: %s",
					id,
					err,
					strings.TrimSpace(string(output)),
				),
			)
		}
	}

	if len(failures) == 0 {
		network := appleNetworkName(projectName)

		exists, err := appleNetworkExists(ctx, network)
		if err != nil {
			failures = append(
				failures,
				fmt.Sprintf(
					"inspect network %s: %v",
					network,
					err,
				),
			)
		} else if exists {
			output, err := appleCommand(
				ctx,
				"network",
				"delete",
				network,
			)

			if err != nil {
				exists, checkErr := appleNetworkExists(ctx, network)

				if checkErr != nil || exists {
					failures = append(
						failures,
						fmt.Sprintf(
							"network %s: %v: %s",
							network,
							err,
							strings.TrimSpace(string(output)),
						),
					)
				}
			}
		}
	}

	if len(failures) > 0 {
		return fmt.Errorf(
			"delete apple container resources: %s",
			strings.Join(failures, "; "),
		)
	}

	// Named volumes are preserved, like Docker Compose down
	// without --volumes.
	fmt.Printf(
		"[apple] deleted project %q\n",
		projectName,
	)

	return nil
}

func appleContainerID(
	project string,
	service string,
) string {
	return project + "-" + service
}

func appleNetworkName(project string) string {
	return "workplace-" + project
}

func appleServiceContainerID(
	project string,
	name string,
	service config.Service,
) string {
	if containerName := strings.TrimSpace(service.ContainerName); containerName != "" {
		return containerName
	}

	return appleContainerID(project, name)
}

func appleUserID() string {
	return strconv.Itoa(os.Getuid())
}

func appleExpandEnvironment(value string) string {
	return strings.ReplaceAll(
		value,
		"${USERID}",
		appleUserID(),
	)
}

func appleServiceNames(cfg config.Config) []string {
	names := make([]string, 0, len(cfg.Services))

	for name, service := range cfg.Services {
		if !service.IsManaged() {
			continue
		}

		// Optional services without an image are disabled.
		if strings.TrimSpace(service.Image) == "" &&
			config.IsOptionalService(name) {
			continue
		}

		names = append(names, name)
	}

	sort.Strings(names)

	return names
}

// appleServiceOrder returns a stable topological ordering.
// Dependencies are started before their dependants.
func appleServiceOrder(
	cfg config.Config,
) ([]string, error) {
	names := appleServiceNames(cfg)

	available := make(map[string]string, len(names))

	for _, name := range names {
		available[name] = name
		available[composeServiceName(name)] = name
	}

	const (
		unvisited = iota
		visiting
		visited
	)

	state := make(map[string]int, len(names))
	var result []string

	var visit func(string) error

	visit = func(name string) error {
		switch state[name] {
		case visiting:
			return fmt.Errorf(
				"circular service dependency involving %q",
				name,
			)
		case visited:
			return nil
		}

		state[name] = visiting

		service := cfg.Services[name]

		for _, dependency := range service.DependsOn {
			dependencyName, exists := available[dependency]

			if !exists {
				dependencyName, exists = available[composeServiceName(dependency)]
			}

			if !exists {
				continue
			}

			if err := visit(dependencyName); err != nil {
				return err
			}
		}

		state[name] = visited
		result = append(result, name)

		return nil
	}

	for _, name := range names {
		if err := visit(name); err != nil {
			return nil, err
		}
	}

	return result, nil
}

func appleNetworkExists(
	ctx context.Context,
	name string,
) (bool, error) {
	networks, err := appleCommandJSON[[]appleNetwork](
		ctx,
		"network",
		"list",
	)

	if err != nil {
		return false, err
	}

	for _, network := range networks {
		if network.ID == name ||
			network.Configuration.Name == name {
			return true, nil
		}
	}

	return false, nil
}

func appleContainerExists(
	ctx context.Context,
	id string,
) bool {
	_, err := appleCommand(
		ctx,
		"inspect",
		id,
	)

	return err == nil
}

func appleContainerStopped(
	ctx context.Context,
	id string,
) bool {
	output, err := appleCommand(
		ctx,
		"inspect",
		id,
	)

	if err != nil {
		return false
	}

	return strings.Contains(
		strings.ToLower(string(output)),
		"stopped",
	)
}

func appleVolumeExists(
	ctx context.Context,
	name string,
) bool {
	_, err := appleCommand(
		ctx,
		"volume",
		"inspect",
		name,
	)

	return err == nil
}

func appleEnsureVolume(
	ctx context.Context,
	name string,
) error {
	if appleVolumeExists(ctx, name) {
		return nil
	}

	output, err := appleCommand(
		ctx,
		"volume",
		"create",
		name,
	)

	if err != nil && !appleVolumeExists(ctx, name) {
		return fmt.Errorf(
			"create volume %q: %w: %s",
			name,
			err,
			strings.TrimSpace(string(output)),
		)
	}

	return nil
}

// appleContainerPort uses Compose's canonical service-to-port mapping.
func appleContainerPort(name string) (string, bool) {
	port, ok := composeContainerPort(name)
	if !ok {
		return "", false
	}

	return strconv.Itoa(port), true
}

func appleMount(
	projectDir string,
	volume config.Volume,
) (string, error) {
	target := strings.TrimSpace(volume.Target)

	if target == "" {
		return "", fmt.Errorf(
			"volume target cannot be empty",
		)
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
			return "", fmt.Errorf(
				"resolve bind mount source: %w",
				err,
			)
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

func applePrepareService(
	ctx context.Context,
	cfg config.Config,
	name string,
	service config.Service,
	network string,
) ([]string, error) {
	if strings.TrimSpace(service.Image) == "" {
		return nil, fmt.Errorf(
			"service %q has no image",
			name,
		)
	}

	args := []string{
		"run",
		"--detach",
		"--name",
		appleServiceContainerID(
			cfg.Name,
			name,
			service,
		),
	}

	// Configure the project network and service aliases.
	composeName := composeServiceName(name)
	networkArgs := network + ",alias=" + name + "." + cfg.Name + ".workplace"

	if composeName != name {
		networkArgs += ",alias=" + composeName + "." + cfg.Name + ".workplace"
	}

	args = append(args, "--network", networkArgs)

	// Keep environment variables deterministic.
	envKeys := make([]string, 0, len(service.Env))

	for key := range service.Env {
		envKeys = append(envKeys, key)
	}

	sort.Strings(envKeys)

	for _, key := range envKeys {
		args = append(
			args,
			"--env",
			key+"="+appleExpandEnvironment(service.Env[key]),
		)
	}

	if service.Port != nil {
		port, ok := composeContainerPort(name)

		if !ok {
			return nil, fmt.Errorf(
				"cannot determine container port for service %q",
				name,
			)
		}

		args = append(
			args,
			"--publish",
			fmt.Sprintf(
				"%d:%d",
				*service.Port,
				port,
			),
		)
	}

	projectDir, err := filepath.Abs(cfg.ProjectDir)
	if err != nil {
		return nil, fmt.Errorf(
			"resolve project directory: %w",
			err,
		)
	}

	for _, volume := range service.Volumes {
		mount, err := appleMount(
			projectDir,
			volume,
		)

		if err != nil {
			return nil, fmt.Errorf(
				"service %q: %w",
				name,
				err,
			)
		}

		if volume.Type == "volume" {
			if err := appleEnsureVolume(
				ctx,
				volume.Name,
			); err != nil {
				return nil, fmt.Errorf(
					"service %q: %w",
					name,
					err,
				)
			}
		}

		args = append(
			args,
			"--volume",
			mount,
		)
	}

	args = append(
		args,
		service.Image,
	)

	return args, nil
}
