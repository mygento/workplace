package cmd

import (
	"fmt"
	"strings"

	"github.com/mygento/workplace/internal/config"
	workplaceruntime "github.com/mygento/workplace/internal/runtime"
	"github.com/spf13/cobra"
)

var composerCmd = &cobra.Command{
	Use:   "composer [arguments...]",
	Short: "Execute Composer inside the PHP container",
	Long:  "Execute Composer commands inside the project's existing PHP container.",

	RunE: func(cmd *cobra.Command, args []string) error {
		return withComposer(cmd, args)
	},
}

func init() {
	addRuntimeFlags(composerCmd)

	// Flags after the first positional argument belong to Composer.
	composerCmd.Flags().SetInterspersed(false)

	rootCmd.AddCommand(composerCmd)
}

func withComposer(
	cmd *cobra.Command,
	args []string,
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

	return workplaceruntime.ExecuteComposer(
		cmd.Context(),
		rt,
		cfg,
		args,
	)
}
