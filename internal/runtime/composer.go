package runtime

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/mygento/workplace/internal/config"
)

// ComposerRunner represents a runtime capable of executing Composer.
type ComposerRunner interface {
	Composer(
		ctx context.Context,
		cfg config.Config,
		args []string,
	) error
}

// dockerPHPContainerName resolves the PHP container name.
//
// Priority:
// 1. Explicit container_name from service configuration.
// 2. Docker Compose default: <project-name>-php-1.
func dockerPHPContainerName(cfg config.Config) (string, error) {
	service, exists := cfg.Services["php"]
	if !exists {
		return "", fmt.Errorf("PHP service is not configured")
	}

	if name := strings.TrimSpace(service.ContainerName); name != "" {
		return name, nil
	}

	if strings.TrimSpace(cfg.Name) == "" {
		return "", fmt.Errorf("project name is not configured")
	}

	return cfg.Name + "-php-1", nil
}

// Composer executes Composer inside an existing PHP container.
// It does not start or recreate containers.
func (r *DockerRuntime) Composer(
	ctx context.Context,
	cfg config.Config,
	args []string,
) error {
	containerName, err := dockerPHPContainerName(cfg)
	if err != nil {
		return err
	}

	if _, err := exec.LookPath("docker"); err != nil {
		return fmt.Errorf("docker CLI not found: %w", err)
	}

	commandArgs := []string{
		"exec",
		"-i",
		containerName,
		"composer",
	}

	commandArgs = append(commandArgs, args...)

	if verbose {
		fmt.Fprintf(
			os.Stderr,
			"+ docker %s\n",
			strings.Join(commandArgs, " "),
		)
	}

	cmd := exec.CommandContext(
		ctx,
		"docker",
		commandArgs...,
	)

	cmd.Env = userCommandEnvironment()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf(
			"execute Composer in container %q: %w",
			containerName,
			err,
		)
	}

	return nil
}

// ExecuteComposer executes Composer using the selected runtime.
func ExecuteComposer(
	ctx context.Context,
	runtime Runtime,
	cfg config.Config,
	args []string,
) error {
	runner, ok := runtime.(ComposerRunner)
	if !ok {
		return fmt.Errorf(
			"Composer execution is not supported by this runtime",
		)
	}

	return runner.Composer(ctx, cfg, args)
}
