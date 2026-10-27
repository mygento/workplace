package runtime

import (
	"context"
	"fmt"
	"os/exec"
	"sort"
	"sync"
	"syscall"

	"github.com/mygento/workplace/internal/config"
)

type NativeRuntime struct {
	mu        sync.Mutex
	processes map[string]*exec.Cmd
}

func NewNative() *NativeRuntime {
	return &NativeRuntime{
		processes: make(map[string]*exec.Cmd),
	}
}

// Start starts managed services using native host executables.
// Unmanaged services are external dependencies and are never started.
func (r *NativeRuntime) Start(
	ctx context.Context,
	cfg config.Config,
) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	names := make([]string, 0, len(cfg.Services))
	for name := range cfg.Services {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		service := cfg.Services[name]

		if !service.IsManaged() {
			fmt.Printf(
				"Skipping unmanaged service %q; expecting an external instance\n",
				name,
			)
			continue
		}

		if _, exists := r.processes[name]; exists {
			continue
		}

		command, args, err := nativeCommand(name)
		if err != nil {
			return err
		}

		path, err := exec.LookPath(command)
		if err != nil {
			return fmt.Errorf(
				"native service %q requires executable %q on PATH: %w",
				name,
				command,
				err,
			)
		}

		cmd := exec.CommandContext(ctx, path, args...)
		cmd.Stdout = nil
		cmd.Stderr = nil

		if err := cmd.Start(); err != nil {
			return fmt.Errorf(
				"start native service %q: %w",
				name,
				err,
			)
		}

		r.processes[name] = cmd

		fmt.Printf(
			"Started native service %q (pid %d)\n",
			name,
			cmd.Process.Pid,
		)

		go func(name string, cmd *exec.Cmd) {
			_ = cmd.Wait()

			r.mu.Lock()
			defer r.mu.Unlock()

			if r.processes[name] == cmd {
				delete(r.processes, name)
			}
		}(name, cmd)
	}

	return nil
}

// Stop terminates only processes started by this runtime.
// It never stops unmanaged external services.
func (r *NativeRuntime) Stop() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	var failures []error

	for name, cmd := range r.processes {
		if cmd.Process == nil {
			delete(r.processes, name)
			continue
		}

		if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
			failures = append(
				failures,
				fmt.Errorf("stop native service %q: %w", name, err),
			)
		}
	}

	if len(failures) > 0 {
		return fmt.Errorf("stop native services: %v", failures)
	}

	return nil
}

// nativeCommand maps supported service names to native executables.
// Container image values are intentionally not used as executable paths.
func nativeCommand(name string) (string, []string, error) {
	switch name {
	case "php":
		return "php-fpm", []string{"-F"}, nil
	case "nginx":
		return "nginx", []string{"-g", "daemon off;"}, nil
	case "mysql":
		return "", nil, fmt.Errorf(
			"native Mysql startup is not implemented; set services.mysql.managed to false in .workplace/config.local.yaml to use an external instance",
		)
	case "opensearch":
		return "", nil, fmt.Errorf(
			"native OpenSearch startup is not implemented; set services.opensearch.managed to false in .workplace/config.local.yaml to use an external instance",
		)
	default:
		return "", nil, fmt.Errorf(
			"native startup is not implemented for service %q",
			name,
		)
	}
}
