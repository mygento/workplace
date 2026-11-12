package cmd

import (
	"fmt"
	"strings"

	"github.com/mygento/workplace/internal/config"
	workplaceruntime "github.com/mygento/workplace/internal/runtime"
	"github.com/spf13/cobra"
)

func withRuntime(
	cmd *cobra.Command,
	action string,
	args ...string,
) error {
	dir, err := cmd.Flags().GetString("dir")
	if err != nil {
		return err
	}

	runtimeName, err := cmd.Flags().GetString("runtime")
	if err != nil {
		return err
	}

	cfg, err := config.NewLoader(dir).Load()
	if err != nil {
		return fmt.Errorf("load project configuration: %w", err)
	}

	var rt workplaceruntime.Runtime

	switch strings.ToLower(strings.TrimSpace(runtimeName)) {
	case "docker":
		rt = workplaceruntime.NewDocker()

	case "apple-container":
		rt = workplaceruntime.NewApple()

	case "native":
		rt = workplaceruntime.NewNative()

	default:
		return fmt.Errorf(
			"unknown runtime %q; supported runtimes: docker, apple-container, native",
			runtimeName,
		)
	}

	switch action {
	case "start":
		return rt.Start(cmd.Context(), cfg)

	case "stop":
		return rt.Stop(cmd.Context(), cfg)

	case "delete":
		return rt.Delete(cmd.Context(), cfg)

	case "logs":
		logsRuntime, ok := rt.(workplaceruntime.LogsRuntime)
		if !ok {
			return fmt.Errorf(
				"logs are not supported by runtime %q",
				runtimeName,
			)
		}

		follow, err := cmd.Flags().GetBool("follow")
		if err != nil {
			return err
		}

		tail, err := cmd.Flags().GetString("tail")
		if err != nil {
			return err
		}

		timestamps, err := cmd.Flags().GetBool("timestamps")
		if err != nil {
			return err
		}

		return logsRuntime.Logs(
			cmd.Context(),
			cfg,
			args,
			follow,
			tail,
			timestamps,
		)

	default:
		return fmt.Errorf("unknown runtime action %q", action)
	}
}

func addRuntimeFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("dir", "d", ".", "Project directory")
	cmd.Flags().String(
		"runtime",
		"docker",
		"Runtime to use: docker, apple-container, or native",
	)
}
