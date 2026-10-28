package cmd

import "github.com/spf13/cobra"

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a workplace environment",
	Long:  "Start the project's managed services using the configured runtime.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return withRuntime(cmd, "start")
	},
}

func init() {
	addRuntimeFlags(startCmd)
	rootCmd.AddCommand(startCmd)
}
