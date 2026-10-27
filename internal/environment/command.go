package environment

import (
	"context"
	"os"
	"os/exec"
	"strings"
)

func run(ctx context.Context, dir, name string, args ...string) commandResult {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = dir

	output, err := cmd.CombinedOutput()
	return commandResult{
		Output: strings.TrimSpace(string(output)),
		Err:    err,
	}
}

func runSimple(name string, args ...string) string {
	output, err := exec.Command(name, args...).CombinedOutput()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func appendUnique(values []string, value string) []string {
	for _, existing := range values {
		if existing == value {
			return values
		}
	}
	return append(values, value)
}
