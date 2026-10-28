package runtime

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/mygento/workplace/internal/config"
)

type AppleRuntime struct{}

func NewApple() *AppleRuntime {
	return &AppleRuntime{}
}

func appleContainerID(project, service string) string {
	return "workplace-" + project + "-" + service
}

func appleNetworkName(project string) string {
	return "workplace-" + project
}

func appleCommand(
	ctx context.Context,
	args ...string,
) ([]byte, error) {
	return RunCommandOutput(ctx, "container", args...)
}

func (r *AppleRuntime) Start(
	ctx context.Context,
	cfg config.Config,
) error {
	if _, err := exec.LookPath("container"); err != nil {
		return fmt.Errorf("apple container CLI not found: %w", err)
	}

	output, err := appleCommand(ctx, "system", "status")
	if err != nil {
		return fmt.Errorf(
			"apple container system is not running; run `container system start`: %w: %s",
			err, strings.TrimSpace(string(output)),
		)
	}

	project := dockerProjectName(cfg)
	network := appleNetworkName(project)

	output, err = appleCommand(ctx, "network", "create", network)
	if err != nil && !appleNetworkExists(ctx, network) {
		return fmt.Errorf(
			"create apple container network %q: %w: %s",
			network, err, strings.TrimSpace(string(output)),
		)
	}

	names := managedAppleServices(cfg)
	var started []string

	rollback := func(cause error) error {
		var failures []string

		for i := len(started) - 1; i >= 0; i-- {
			id := started[i]
			output, err := appleCommand(
				context.WithoutCancel(ctx),
				"delete", "--force", id,
			)
			if err != nil && appleContainerExists(
				context.WithoutCancel(ctx), id,
			) {
				failures = append(failures, fmt.Sprintf(
					"%s: %v: %s",
					id, err, strings.TrimSpace(string(output)),
				))
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
		svc := cfg.Services[name]
		if strings.TrimSpace(svc.Image) == "" {
			return rollback(fmt.Errorf("service %q has no image", name))
		}

		id := appleContainerID(project, name)
		if appleContainerExists(ctx, id) {
			fmt.Printf(
				"apple container service %q already exists; skipping\n",
				name,
			)
			continue
		}

		args := []string{
			"run",
			"--detach",
			"--name", id,
			"--network", network + ",alias=" + name,
		}

		envKeys := make([]string, 0, len(svc.Env))
		for key := range svc.Env {
			envKeys = append(envKeys, key)
		}
		sort.Strings(envKeys)

		for _, key := range envKeys {
			args = append(args, "--env", key+"="+svc.Env[key])
		}

		if svc.Port != nil {
			port, ok := appleContainerPort(name)
			if !ok {
				return rollback(fmt.Errorf(
					"unknown container port for service %q",
					name,
				))
			}
			args = append(
				args,
				"--publish",
				fmt.Sprintf("%d:%s", *svc.Port, port),
			)
		}

		if cfg.ProjectDir != "" {
			source, err := filepath.Abs(cfg.ProjectDir)
			if err != nil {
				return rollback(fmt.Errorf("resolve project directory: %w", err))
			}
			args = append(args, "--volume", source+":/workspace")
		}

		args = append(args, svc.Image)

		output, err := appleCommand(ctx, args...)
		if err != nil {
			return rollback(fmt.Errorf(
				"start apple container service %q: %w: %s",
				name, err, strings.TrimSpace(string(output)),
			))
		}

		started = append(started, id)
		fmt.Printf("started apple container service %q\n", name)
	}

	return nil
}

func (r *AppleRuntime) Stop(
	ctx context.Context,
	cfg config.Config,
) error {
	if _, err := exec.LookPath("container"); err != nil {
		return fmt.Errorf("apple container CLI not found: %w", err)
	}

	project := dockerProjectName(cfg)
	var failures []string

	for _, name := range managedAppleServices(cfg) {
		id := appleContainerID(project, name)
		if !appleContainerExists(ctx, id) {
			continue
		}

		stopCtx, cancel := context.WithTimeout(ctx, 20*time.Second)
		output, err := appleCommand(
			stopCtx, "stop", "--time", "10", id,
		)
		cancel()

		if err != nil && !appleContainerStopped(ctx, id) {
			failures = append(failures, fmt.Sprintf(
				"%s: %v: %s",
				id, err, strings.TrimSpace(string(output)),
			))
		}
	}

	if len(failures) > 0 {
		return fmt.Errorf(
			"stop apple container services: %s",
			strings.Join(failures, "; "),
		)
	}

	return nil
}

func (r *AppleRuntime) Delete(
	ctx context.Context,
	cfg config.Config,
) error {
	if _, err := exec.LookPath("container"); err != nil {
		return fmt.Errorf("apple container CLI not found: %w", err)
	}

	project := dockerProjectName(cfg)
	var failures []string

	for _, name := range managedAppleServices(cfg) {
		id := appleContainerID(project, name)
		if !appleContainerExists(ctx, id) {
			continue
		}

		output, err := appleCommand(
			ctx, "delete", "--force", id,
		)
		if err != nil && appleContainerExists(ctx, id) {
			failures = append(failures, fmt.Sprintf(
				"%s: %v: %s",
				id, err, strings.TrimSpace(string(output)),
			))
		}
	}

	if len(failures) == 0 {
		network := appleNetworkName(project)
		output, err := appleCommand(ctx, "network", "delete", network)
		if err != nil && appleNetworkExists(ctx, network) {
			failures = append(failures, fmt.Sprintf(
				"network %s: %v: %s",
				network, err, strings.TrimSpace(string(output)),
			))
		}
	}

	if len(failures) > 0 {
		return fmt.Errorf(
			"delete apple container resources: %s",
			strings.Join(failures, "; "),
		)
	}

	return nil
}

func managedAppleServices(cfg config.Config) []string {
	names := make([]string, 0, len(cfg.Services))

	for name, svc := range cfg.Services {
		if svc.IsManaged() {
			names = append(names, name)
		}
	}

	sort.Strings(names)
	return names
}

func appleContainerExists(ctx context.Context, id string) bool {
	_, err := appleCommand(ctx, "inspect", id)
	return err == nil
}

func appleContainerStopped(ctx context.Context, id string) bool {
	output, err := appleCommand(ctx, "inspect", id)
	if err != nil {
		return false
	}

	return strings.Contains(
		strings.ToLower(string(output)),
		"stopped",
	)
}

func appleNetworkExists(ctx context.Context, network string) bool {
	output, err := appleCommand(ctx, "network", "list")
	if err != nil {
		return false
	}

	return strings.Contains(string(output), network)
}

func appleContainerPort(name string) (string, bool) {
	port, ok := map[string]string{
		"nginx":      "80",
		"php":        "9000",
		"mysql":      "3306",
		"opensearch": "9200",
	}[name]

	return port, ok
}
