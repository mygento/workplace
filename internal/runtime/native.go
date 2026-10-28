package runtime

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/mygento/workplace/internal/config"
)

type NativeRuntime struct{}

type nativeProcess struct {
	PID         int    `json:"pid"`
	Executable  string `json:"executable"`
	ServiceName string `json:"service"`
}

func NewNative() *NativeRuntime {
	return &NativeRuntime{}
}

func nativeRuntimeDir(cfg config.Config) string {
	return filepath.Join(cfg.ProjectDir, ".workplace", "runtime", "native")
}

func nativeStatePath(cfg config.Config, name string) string {
	return filepath.Join(nativeRuntimeDir(cfg), name+".json")
}

func (r *NativeRuntime) Start(
	ctx context.Context,
	cfg config.Config,
) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if err := os.MkdirAll(nativeRuntimeDir(cfg), 0700); err != nil {
		return fmt.Errorf("create native runtime directory: %w", err)
	}

	names := make([]string, 0, len(cfg.Services))
	for name := range cfg.Services {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		svc := cfg.Services[name]
		if !svc.IsManaged() {
			fmt.Printf("skipping unmanaged service %q\n", name)
			continue
		}

		statePath := nativeStatePath(cfg, name)
		state, err := readNativeState(statePath)
		if err == nil {
			if nativeProcessMatches(ctx, state) {
				fmt.Printf(
					"native service %q is already running (pid %d)\n",
					name, state.PID,
				)
				continue
			}

			if err := os.Remove(statePath); err != nil &&
				!errors.Is(err, os.ErrNotExist) {
				return fmt.Errorf(
					"remove stale state for native service %q: %w",
					name, err,
				)
			}
		} else if !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf(
				"read state for native service %q: %w",
				name, err,
			)
		}

		command, args, err := nativeCommand(name)
		if err != nil {
			return err
		}

		executable, err := exec.LookPath(command)
		if err != nil {
			return fmt.Errorf(
				"native service %q requires executable %q on PATH: %w",
				name, command, err,
			)
		}

		cmd := exec.Command(executable, args...)
		cmd.Stdout = nil
		cmd.Stderr = nil
		cmd.Stdin = nil

		if verbose {
			fmt.Fprintf(os.Stderr, "+ %s\n", shellCommand(executable, args))
		}

		if err := cmd.Start(); err != nil {
			return fmt.Errorf("start native service %q: %w", name, err)
		}

		state = nativeProcess{
			PID:         cmd.Process.Pid,
			Executable:  executable,
			ServiceName: name,
		}

		if err := writeNativeState(statePath, state); err != nil {
			_ = cmd.Process.Signal(syscall.SIGTERM)
			_ = cmd.Wait()
			return fmt.Errorf(
				"persist state for native service %q: %w",
				name, err,
			)
		}

		go func() {
			_ = cmd.Wait()
		}()

		fmt.Printf(
			"started native service %q (pid %d)\n",
			name, state.PID,
		)
	}

	return nil
}

func (r *NativeRuntime) Stop(
	ctx context.Context,
	cfg config.Config,
) error {
	return stopNativeServices(ctx, cfg, false)
}

func (r *NativeRuntime) Delete(
	ctx context.Context,
	cfg config.Config,
) error {
	if err := stopNativeServices(ctx, cfg, true); err != nil {
		return err
	}

	err := os.Remove(nativeRuntimeDir(cfg))
	if err != nil &&
		!errors.Is(err, os.ErrNotExist) &&
		!errors.Is(err, syscall.ENOTEMPTY) {
		return fmt.Errorf("remove native runtime directory: %w", err)
	}

	return nil
}

func stopNativeServices(
	ctx context.Context,
	cfg config.Config,
	removeState bool,
) error {
	names := make([]string, 0, len(cfg.Services))
	for name, svc := range cfg.Services {
		if svc.IsManaged() {
			names = append(names, name)
		}
	}
	sort.Strings(names)

	var failures []string

	for _, name := range names {
		if err := ctx.Err(); err != nil {
			return err
		}

		statePath := nativeStatePath(cfg, name)
		state, err := readNativeState(statePath)
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil {
			failures = append(failures, fmt.Sprintf(
				"%s: read process state: %v", name, err,
			))
			continue
		}

		if !nativeProcessMatches(ctx, state) {
			if err := os.Remove(statePath); err != nil &&
				!errors.Is(err, os.ErrNotExist) {
				failures = append(failures, fmt.Sprintf(
					"%s: remove stale state: %v", name, err,
				))
			}
			continue
		}

		process, err := os.FindProcess(state.PID)
		if err != nil {
			failures = append(failures, fmt.Sprintf(
				"%s: find process: %v", name, err,
			))
			continue
		}

		if verbose {
			fmt.Fprintf(
				os.Stderr,
				"+ kill -TERM %d\n",
				state.PID,
			)
		}

		if err := process.Signal(syscall.SIGTERM); err != nil &&
			!errors.Is(err, os.ErrProcessDone) {
			failures = append(failures, fmt.Sprintf(
				"%s: send SIGTERM: %v", name, err,
			))
			continue
		}

		deadline := time.Now().Add(10 * time.Second)
		stopped := false

		for time.Now().Before(deadline) {
			if err := ctx.Err(); err != nil {
				failures = append(failures, fmt.Sprintf(
					"%s: interrupted while waiting: %v", name, err,
				))
				break
			}

			if !nativeProcessMatches(ctx, state) {
				stopped = true
				break
			}

			time.Sleep(100 * time.Millisecond)
		}

		if !stopped && nativeProcessMatches(ctx, state) {
			failures = append(failures, fmt.Sprintf(
				"%s: process %d did not exit after SIGTERM",
				name, state.PID,
			))
			continue
		}

		if removeState {
			if err := os.Remove(statePath); err != nil &&
				!errors.Is(err, os.ErrNotExist) {
				failures = append(failures, fmt.Sprintf(
					"%s: remove state: %v", name, err,
				))
				continue
			}
		}

		fmt.Printf("stopped native service %q\n", name)
	}

	if len(failures) > 0 {
		return fmt.Errorf(
			"native service lifecycle errors: %s",
			strings.Join(failures, "; "),
		)
	}

	return nil
}

func readNativeState(path string) (nativeProcess, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nativeProcess{}, err
	}

	var state nativeProcess
	if err := json.Unmarshal(data, &state); err != nil {
		return nativeProcess{}, fmt.Errorf("parse state: %w", err)
	}

	if state.PID <= 0 ||
		state.Executable == "" ||
		state.ServiceName == "" {
		return nativeProcess{}, fmt.Errorf(
			"invalid process state in %s", path,
		)
	}

	return state, nil
}

func writeNativeState(path string, state nativeProcess) error {
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}

	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0600); err != nil {
		return err
	}

	if err := os.Rename(tmp, path); err != nil {
		_ = os.Remove(tmp)
		return err
	}

	return nil
}

func nativeProcessMatches(
	ctx context.Context,
	state nativeProcess,
) bool {
	if state.PID <= 0 {
		return false
	}

	output, err := RunCommandOutput(
		ctx,
		"ps",
		"-p", strconv.Itoa(state.PID),
		"-o", "command=",
	)
	if err != nil {
		return false
	}

	commandLine := strings.TrimSpace(string(output))
	if commandLine == "" {
		return false
	}

	return strings.Contains(commandLine, state.Executable)
}

func nativeCommand(name string) (string, []string, error) {
	switch name {
	case "php":
		return "php-fpm", []string{"-F"}, nil
	case "nginx":
		return "nginx", []string{"-g", "daemon off;"}, nil
	case "mysql":
		return "", nil, fmt.Errorf(
			"native MySQL startup is not implemented; configure an external instance or disable managed startup",
		)
	case "opensearch":
		return "", nil, fmt.Errorf(
			"native OpenSearch startup is not implemented; configure an external instance or disable managed startup",
		)
	default:
		return "", nil, fmt.Errorf(
			"native startup is not implemented for service %q", name,
		)
	}
}
