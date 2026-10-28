package runtime

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

var verbose bool

func SetVerbose(enabled bool) {
	verbose = enabled
}

func runCommandWithEnv(
	ctx context.Context,
	env []string,
	name string,
	args ...string,
) ([]byte, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Env = env

	if verbose {
		fmt.Fprintf(os.Stderr, "+ %s\n", shellCommand(name, args))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return nil, fmt.Errorf("run command %q: %w", name, err)
		}

		return nil, nil
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return output, fmt.Errorf("run command %q: %w", name, err)
	}

	return output, nil
}

func RunCommand(
	ctx context.Context,
	name string,
	args ...string,
) error {
	_, err := RunCommandOutput(ctx, name, args...)
	return err
}

func RunCommandOutput(
	ctx context.Context,
	name string,
	args ...string,
) ([]byte, error) {
	cmd := exec.CommandContext(ctx, name, args...)

	if verbose {
		fmt.Fprintf(os.Stderr, "+ %s\n", shellCommand(name, args))

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return nil, fmt.Errorf("run command %q: %w", name, err)
		}

		return nil, nil
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return output, fmt.Errorf("run command %q: %w", name, err)
	}

	return output, nil
}

func RunCommandWithOutput(
	ctx context.Context,
	name string,
	stdout io.Writer,
	stderr io.Writer,
	args ...string,
) error {
	if verbose {
		fmt.Fprintf(os.Stderr, "+ %s\n", shellCommand(name, args))
	}

	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("run command %q: %w", name, err)
	}

	return nil
}

func shellCommand(name string, args []string) string {
	parts := make([]string, 0, len(args)+1)
	parts = append(parts, shellQuote(name))

	for _, arg := range args {
		parts = append(parts, shellQuote(arg))
	}

	return strings.Join(parts, " ")
}

func shellQuote(value string) string {
	if value == "" {
		return "''"
	}

	if strings.ContainsAny(value, " \t\n'\"\\$`!&;|()<>{}*?[]") {
		return "'" + strings.ReplaceAll(value, "'", "'\\''") + "'"
	}

	return value
}

// RunCommandOutputCapture captures stdout and stderr even in verbose mode.
// Use this for commands whose output must be parsed by the caller.
func RunCommandOutputCapture(
	ctx context.Context,
	name string,
	args ...string,
) ([]byte, error) {
	if verbose {
		fmt.Fprintf(os.Stderr, "+ %s\n", shellCommand(name, args))
	}

	var output bytes.Buffer

	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Stdout = &output
	cmd.Stderr = &output

	if err := cmd.Run(); err != nil {
		return output.Bytes(), fmt.Errorf("run command %q: %w", name, err)
	}

	return output.Bytes(), nil
}
