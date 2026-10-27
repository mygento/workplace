package environment

import (
	"context"
	"os/exec"
)

func collectNode(ctx context.Context, dir string) NodeInfo {
	return NodeInfo{
		Node: inspectTool(ctx, dir, "node", "--version"),
		NPM:  inspectTool(ctx, dir, "npm", "--version"),
	}
}

func inspectTool(
	ctx context.Context,
	dir string,
	name string,
	args ...string,
) ToolInfo {
	var result ToolInfo

	path, err := exec.LookPath(name)
	if err != nil {
		result.Error = name + " executable not found in PATH"
		return result
	}

	result.Path = path

	version := run(ctx, dir, path, args...)
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
