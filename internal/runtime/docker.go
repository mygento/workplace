package runtime

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/mygento/workplace/internal/config"
)

type DockerRuntime struct{}

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

	if name == "" || name == "." {
		name = "workplace"
	}

	if name[0] < 'a' || name[0] > 'z' {
		name = "workplace-" + name
	}

	return name
}

func dockerRuntimeDir(cfg config.Config) string {
	return filepath.Join(cfg.ProjectDir, ".workplace", "runtime")
}

func dockerComposePath(cfg config.Config) string {
	return filepath.Join(dockerRuntimeDir(cfg), "compose.yaml")
}

func dockerCommandEnvironment(cfg config.Config) []string {
	env := os.Environ()

	env = setEnvironmentValue(
		env,
		"COMPOSE_PROJECT_NAME",
		dockerProjectName(cfg),
	)

	env = setEnvironmentValue(
		env,
		"USERID",
		strconv.Itoa(os.Getuid()),
	)

	return env
}

func setEnvironmentValue(env []string, key, value string) []string {
	prefix := key + "="
	result := make([]string, 0, len(env)+1)

	for _, entry := range env {
		if !strings.HasPrefix(entry, prefix) {
			result = append(result, entry)
		}
	}

	return append(result, prefix+value)
}

func (r *DockerRuntime) runCompose(
	ctx context.Context,
	cfg config.Config,
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
		"--project-name", dockerProjectName(cfg),
		"--file", composePath,
	}
	commandArgs = append(commandArgs, args...)

	return runCommandWithEnv(
		ctx,
		dockerCommandEnvironment(cfg),
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
		return fmt.Errorf("render Docker Compose configuration: %w", err)
	}

	output, err := r.runCompose(
		ctx,
		cfg,
		"up",
		"-d",
		"--remove-orphans",
	)
	if err != nil {
		return fmt.Errorf(
			"start Docker Compose project %q: %w: %s",
			dockerProjectName(cfg),
			err,
			strings.TrimSpace(string(output)),
		)
	}

	fmt.Printf(
		"started Docker Compose project %q\n",
		dockerProjectName(cfg),
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

	output, err := r.runCompose(ctx, cfg, "stop")
	if err != nil {
		return fmt.Errorf(
			"stop Docker Compose project %q: %w: %s",
			dockerProjectName(cfg),
			err,
			strings.TrimSpace(string(output)),
		)
	}

	fmt.Printf(
		"stopped Docker Compose project %q\n",
		dockerProjectName(cfg),
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

	output, err := r.runCompose(ctx, cfg, "down")
	if err != nil {
		return fmt.Errorf(
			"delete Docker Compose project %q: %w: %s",
			dockerProjectName(cfg),
			err,
			strings.TrimSpace(string(output)),
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
		dockerProjectName(cfg),
	)
	return nil
}
