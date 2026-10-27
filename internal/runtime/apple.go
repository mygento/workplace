package runtime

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/mygento/workplace/internal/config"
)

type AppleRuntime struct {
	mu         sync.Mutex
	project    string
	network    string
	containers []string
}

func NewApple() *AppleRuntime {
	return &AppleRuntime{}
}

func (r *AppleRuntime) Start(
	ctx context.Context,
	cfg config.Config,
) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, err := exec.LookPath("container"); err != nil {
		return fmt.Errorf("apple container CLI not found: %w", err)
	}

	project := dockerProjectName(cfg)
	network := "workplace-" + project

	if output, err := exec.CommandContext(
		ctx, "container", "system", "status",
	).CombinedOutput(); err != nil {
		return fmt.Errorf(
			"apple container system is not running; run `container system start`: %w: %s",
			err, strings.TrimSpace(string(output)),
		)
	}

	// On macOS 26+, user-defined networks are supported.
	if output, err := exec.CommandContext(
		ctx, "container", "network", "create", network,
	).CombinedOutput(); err != nil {
		return fmt.Errorf(
			"create Apple Container network %q: %w: %s",
			network, err, strings.TrimSpace(string(output)),
		)
	}
	r.project = project
	r.network = network

	names := make([]string, 0, len(cfg.Services))
	for name := range cfg.Services {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		svc := cfg.Services[name]
		if !svc.IsManaged() {
			continue
		}
		if strings.TrimSpace(svc.Image) == "" {
			return r.rollback(ctx, fmt.Errorf(
				"service %q has no image", name,
			))
		}

		id := "workplace-" + project + "-" + name
		args := []string{
			"run", "--detach",
			"--name", id,
			"--network", network + ",alias=" + name,
		}

		keys := make([]string, 0, len(svc.Env))
		for key := range svc.Env {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			args = append(args, "--env", key+"="+svc.Env[key])
		}

		if svc.Port != nil {
			port, ok := appleContainerPort(name)
			if !ok {
				return r.rollback(ctx, fmt.Errorf(
					"unknown container port for service %q", name,
				))
			}
			args = append(args, "--publish",
				fmt.Sprintf("%d:%s", *svc.Port, port))
		}

		if cfg.ProjectDir != "" {
			source, err := filepath.Abs(cfg.ProjectDir)
			if err != nil {
				return r.rollback(ctx, err)
			}
			args = append(args, "--volume", source+":/workspace")
		}

		args = append(args, svc.Image)
		cmd := exec.CommandContext(ctx, "container", args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return r.rollback(ctx, fmt.Errorf(
				"start Apple Container service %q: %w", name, err,
			))
		}

		r.containers = append(r.containers, id)
		fmt.Printf("Started Apple Container service %q\n", name)
	}

	return nil
}

func (r *AppleRuntime) Stop() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	var failures []string

	for i := len(r.containers) - 1; i >= 0; i-- {
		ctx, cancel := context.WithTimeout(
			context.Background(), 20*time.Second,
		)
		output, err := exec.CommandContext(
			ctx, "container", "stop", "--time", "10",
			r.containers[i],
		).CombinedOutput()
		cancel()

		if err != nil {
			failures = append(failures, fmt.Sprintf(
				"%s: %v: %s",
				r.containers[i], err, strings.TrimSpace(string(output)),
			))
		}
	}

	if len(failures) == 0 && r.network != "" {
		output, err := exec.Command(
			"container", "network", "delete", r.network,
		).CombinedOutput()
		if err != nil {
			failures = append(failures, fmt.Sprintf(
				"network %s: %v: %s",
				r.network, err, strings.TrimSpace(string(output)),
			))
		}
	}

	if len(failures) > 0 {
		return fmt.Errorf(
			"apple container cleanup failed: %s",
			strings.Join(failures, "; "),
		)
	}

	r.containers = nil
	r.network = ""
	r.project = ""
	return nil
}

func (r *AppleRuntime) rollback(
	ctx context.Context,
	cause error,
) error {
	var failures []string

	for i := len(r.containers) - 1; i >= 0; i-- {
		out, err := exec.CommandContext(
			ctx, "container", "stop", "--time", "2",
			r.containers[i],
		).CombinedOutput()
		if err != nil {
			failures = append(failures, fmt.Sprintf(
				"%s: %v: %s",
				r.containers[i], err, strings.TrimSpace(string(out)),
			))
		}
	}

	if r.network != "" {
		out, err := exec.CommandContext(
			ctx, "container", "network", "delete", r.network,
		).CombinedOutput()
		if err != nil {
			failures = append(failures, fmt.Sprintf(
				"network %s: %v: %s",
				r.network, err, strings.TrimSpace(string(out)),
			))
		}
	}

	r.containers = nil
	r.network = ""
	r.project = ""

	if len(failures) > 0 {
		return fmt.Errorf("%w; rollback errors: %s",
			cause, strings.Join(failures, "; "))
	}
	return cause
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
