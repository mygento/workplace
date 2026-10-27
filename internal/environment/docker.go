package environment

import (
	"context"
	"os/exec"
	"strings"
)

func collectDocker(
	ctx context.Context,
	dir string,
	containers *[]ContainerInfo,
) DockerInfo {
	var result DockerInfo

	path, err := exec.LookPath("docker")
	if err != nil {
		return result
	}

	result.ClientInstalled = true
	result.ClientPath = path

	client := run(
		ctx, dir, path,
		"version", "--format", "{{.Client.Version}}",
	)
	if client.Err == nil {
		result.ClientVersion = client.Output
	} else {
		result.ClientVersion = runSimple(path, "--version")
	}

	server := run(
		ctx, dir, path,
		"version", "--format", "{{.Server.Version}}",
	)
	if server.Err == nil && server.Output != "" {
		result.DaemonAvailable = true
		result.DaemonVersion = server.Output
	} else {
		result.DaemonError = server.Output
		if result.DaemonError == "" && server.Err != nil {
			result.DaemonError = server.Err.Error()
		}
		if result.DaemonError == "" {
			result.DaemonError = "Docker server version unavailable"
		}
	}

	compose := run(ctx, dir, path, "compose", "version", "--short")
	if compose.Err == nil {
		result.ComposeInstalled = true
		result.ComposeVersion = compose.Output
		result.ComposeCommand = []string{path, "compose"}
	} else {
		legacyPath, lookupErr := exec.LookPath("docker-compose")
		if lookupErr != nil {
			result.ComposeError = compose.Output
			if result.ComposeError == "" {
				result.ComposeError = "Docker Compose not found"
			}
		} else {
			legacy := run(ctx, dir, legacyPath, "version", "--short")
			if legacy.Err == nil {
				result.ComposeInstalled = true
				result.ComposeVersion = legacy.Output
				result.ComposeCommand = []string{legacyPath}
			} else {
				result.ComposeError = legacy.Output
				if result.ComposeError == "" {
					result.ComposeError = legacy.Err.Error()
				}
			}
		}
	}

	if result.DaemonAvailable {
		collectContainers(ctx, dir, path, containers)
	}

	return result
}

func collectContainers(
	ctx context.Context,
	dir string,
	dockerPath string,
	containers *[]ContainerInfo,
) {
	result := run(
		ctx, dir, dockerPath,
		"ps",
		"--format",
		"{{.ID}}\t{{.Names}}\t{{.Image}}\t{{.Status}}",
	)
	if result.Err != nil {
		return
	}

	for _, line := range strings.Split(result.Output, "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}

		fields := strings.Split(line, "\t")
		if len(fields) < 4 {
			continue
		}

		*containers = append(*containers, ContainerInfo{
			ID:     fields[0],
			Name:   fields[1],
			Image:  fields[2],
			Status: fields[3],
		})
	}
}
