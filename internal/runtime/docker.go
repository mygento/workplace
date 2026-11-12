package runtime

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/mygento/workplace/internal/config"
	"github.com/mygento/workplace/internal/project"
)

type DockerRuntime struct{}

func NewDocker() *DockerRuntime {
	return &DockerRuntime{}
}

func dockerRuntimeDir(cfg config.Config) string {
	return filepath.Join(cfg.ProjectDir, ".workplace", "runtime")
}

func dockerComposePath(cfg config.Config) string {
	return filepath.Join(dockerRuntimeDir(cfg), "compose.yaml")
}

func (r *DockerRuntime) runCompose(
	ctx context.Context,
	cfg config.Config,
	interactive bool,
	args ...string,
) ([]byte, error) {
	if _, err := exec.LookPath("docker"); err != nil {
		return nil, fmt.Errorf("docker CLI not found: %w", err)
	}

	composePath := dockerComposePath(cfg)

	if _, err := os.Stat(composePath); err != nil {
		return nil, fmt.Errorf(
			"access Compose configuration %s: %w",
			composePath,
			err,
		)
	}

	commandArgs := []string{
		"compose",
		"--project-name", cfg.Name,
		"--file", composePath,
	}

	commandArgs = append(commandArgs, args...)

	env := userCommandEnvironment()

	if interactive {
		err := RunInteractiveCommand(
			ctx,
			env,
			"docker",
			commandArgs...,
		)

		return nil, err
	}

	return RunCommand(
		ctx,
		env,
		"docker",
		commandArgs...,
	)
}

func (r *DockerRuntime) Start(
	ctx context.Context,
	cfg config.Config,
) error {
	if cfg.ProjectDir == "" {
		return fmt.Errorf("project directory is not configured")
	}

	composePath := dockerComposePath(cfg)

	if err := writeCompose(cfg, composePath); err != nil {
		return fmt.Errorf(
			"render Docker Compose configuration: %w",
			err,
		)
	}

	if cfg.HasService("nginx") {
		if err := project.Install(cfg.ProjectDir, cfg.Type); err != nil {
			return fmt.Errorf(
				"install project nginx configuration: %w",
				err,
			)
		}
	}

	_, err := r.runCompose(
		ctx,
		cfg,
		true,
		"up",
		"-d",
		"--remove-orphans",
	)

	if err != nil {
		return fmt.Errorf(
			"start Docker Compose project %q: %w",
			cfg.Name,
			err,
		)
	}

	fmt.Printf(
		"started Docker Compose project %q\n",
		cfg.Name,
	)

	return nil
}

func (r *DockerRuntime) Stop(
	ctx context.Context,
	cfg config.Config,
) error {
	composePath := dockerComposePath(cfg)

	if _, err := os.Stat(composePath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}

		return fmt.Errorf("access Compose configuration: %w", err)
	}

	_, err := r.runCompose(
		ctx,
		cfg,
		true,
		"stop",
	)

	if err != nil {
		return fmt.Errorf(
			"stop Docker Compose project %q: %w",
			cfg.Name,
			err,
		)
	}

	fmt.Printf(
		"stopped Docker Compose project %q\n",
		cfg.Name,
	)

	return nil
}

func (r *DockerRuntime) Delete(
	ctx context.Context,
	cfg config.Config,
) error {
	composePath := dockerComposePath(cfg)

	if _, err := os.Stat(composePath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}

		return fmt.Errorf("access Compose configuration: %w", err)
	}

	_, err := r.runCompose(
		ctx,
		cfg,
		true,
		"down",
	)

	if err != nil {
		return fmt.Errorf(
			"delete Docker Compose project %q: %w",
			cfg.Name,
			err,
		)
	}

	if err := os.Remove(composePath); err != nil &&
		!errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf(
			"remove generated Compose configuration: %w",
			err,
		)
	}

	fmt.Printf(
		"deleted Docker Compose project %q\n",
		cfg.Name,
	)

	return nil
}

// Logs streams Docker Compose service logs directly to the terminal.
func (r *DockerRuntime) Logs(
	ctx context.Context,
	cfg config.Config,
	services []string,
	follow bool,
	tail string,
	timestamps bool,
) error {
	args := []string{"logs"}

	if follow {
		args = append(args, "--follow")
	}

	if tail != "" {
		args = append(args, "--tail", tail)
	}

	if timestamps {
		args = append(args, "--timestamps")
	}

	args = append(args, services...)

	_, err := r.runCompose(
		ctx,
		cfg,
		true,
		args...,
	)

	if err != nil {
		return fmt.Errorf(
			"show Docker Compose logs for project %q: %w",
			cfg.Name,
			err,
		)
	}

	return nil
}
