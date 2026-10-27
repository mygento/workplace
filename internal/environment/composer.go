package environment

import (
	"context"
	"os/exec"
)

func collectComposer(ctx context.Context, dir string) ToolInfo {
	var result ToolInfo

	path, err := exec.LookPath("composer")
	if err != nil {
		result.Error = "composer executable not found in PATH"
		return result
	}

	result.Path = path

	version := run(ctx, dir, path, "--version", "--no-ansi")
	if version.Err != nil {
		result.Error = version.Output
		if result.Error == "" {
			result.Error = version.Err.Error()
		}
		return result
	}

	result.Installed = true
	result.Version = version.Output
	return result
}
