package cmd

import "github.com/spf13/cobra"

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop a workplace environment",
	Long:  "Stop the project's managed services without deleting their resources.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return withRuntime(cmd, "stop")
	},
}

func init() {
	addRuntimeFlags(stopCmd)
	rootCmd.AddCommand(stopCmd)
}
