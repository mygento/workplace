package cmd

import (
	"fmt"
	"strings"

	"github.com/mygento/workplace/internal/config"
	workplaceruntime "github.com/mygento/workplace/internal/runtime"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a workplace environment",
	Long: `Load the project configuration and start its services
using Docker, Apple container, or native processes.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, err := cmd.Flags().GetString("dir")
		if err != nil {
			return err
		}

		runtimeName, err := cmd.Flags().GetString("runtime")
		if err != nil {
			return err
		}

		switch strings.ToLower(strings.TrimSpace(runtimeName)) {
		case "docker", "apple-container", "native":
			cfg, err := config.NewLoader(dir).Load()
			if err != nil {
				return fmt.Errorf(
					"load project configuration: %w",
					err,
				)
			}

			switch strings.ToLower(strings.TrimSpace(runtimeName)) {
			case "docker":
				runtime := workplaceruntime.NewDocker()
				return runtime.Start(cmd.Context(), cfg)

			case "apple-container":
				runtime := workplaceruntime.NewApple()
				return runtime.Start(cmd.Context(), cfg)

			case "native":
				runtime := workplaceruntime.NewNative()
				return runtime.Start(cmd.Context(), cfg)
			}
		}

		return fmt.Errorf(
			"unknown runtime %q; supported runtimes: docker, apple-container, native",
			runtimeName,
		)
	},
}

func init() {
	startCmd.Flags().StringP(
		"dir",
		"d",
		".",
		"Project directory",
	)

	startCmd.Flags().String(
		"runtime",
		"docker",
		"Runtime to use: docker, apple-container, or native",
	)

	rootCmd.AddCommand(startCmd)
}
