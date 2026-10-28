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

	var runtime workplaceruntime.Runtime

	switch strings.ToLower(strings.TrimSpace(runtimeName)) {
	case "docker":
		runtime = workplaceruntime.NewDocker()
	case "apple-container":
		runtime = workplaceruntime.NewApple()
	case "native":
		runtime = workplaceruntime.NewNative()
	default:
		return fmt.Errorf(
			"unknown runtime %q; supported runtimes: docker, apple-container, native",
			runtimeName,
		)
	}

	switch action {
	case "start":
		return runtime.Start(cmd.Context(), cfg)
	case "stop":
		return runtime.Stop(cmd.Context(), cfg)
	case "delete":
		return runtime.Delete(cmd.Context(), cfg)
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
