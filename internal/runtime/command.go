package runtime

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var verbose bool

func SetVerbose(enabled bool) {
	verbose = enabled
}

func userCommandEnvironment() []string {
	return setEnvironmentValue(
		os.Environ(),
		"USERID",
		strconv.Itoa(os.Getuid()),
	)
}

func setEnvironmentValue(
	env []string,
	key string,
	value string,
) []string {
	prefix := key + "="
	result := make([]string, 0, len(env)+1)

	for _, entry := range env {
		if !strings.HasPrefix(entry, prefix) {
			result = append(result, entry)
		}
	}

	return append(result, prefix+value)
}

func RunCommand(
	ctx context.Context,
	env []string,
	name string,
	args ...string,
) ([]byte, error) {
	cmd := exec.CommandContext(ctx, name, args...)

	if env != nil {
		cmd.Env = env
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	if verbose {
		fmt.Fprintf(os.Stderr, "+ %s\n", shellCommand(name, args))
		cmd.Stdout = io.MultiWriter(os.Stdout, &stdout)
		cmd.Stderr = io.MultiWriter(os.Stderr, &stderr)
	} else {
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
	}

	err := cmd.Run()

	output := append(stdout.Bytes(), stderr.Bytes()...)

	if err != nil {
		return output, fmt.Errorf(
			"run command %q: %w",
			name,
			err,
		)
	}

	return output, nil
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
