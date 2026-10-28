package cmd

import "github.com/spf13/cobra"

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a workplace environment",
	Long:  "Remove the project's managed runtime resources.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return withRuntime(cmd, "delete")
	},
}

func init() {
	addRuntimeFlags(deleteCmd)
	rootCmd.AddCommand(deleteCmd)
}
