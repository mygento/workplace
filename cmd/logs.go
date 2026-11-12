package cmd

import "github.com/spf13/cobra"

var (
	logsFollow     bool
	logsTail       string
	logsTimestamps bool

	logsCmd = &cobra.Command{
		Use:   "logs [SERVICE...]",
		Short: "Show service logs",
		Long:  "Show logs from the project's managed services.",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return withRuntime(cmd, "logs", args...)
		},
	}
)

func init() {
	addRuntimeFlags(logsCmd)

	logsCmd.Flags().BoolVarP(
		&logsFollow,
		"follow",
		"f",
		false,
		"Follow log output",
	)

	logsCmd.Flags().StringVar(
		&logsTail,
		"tail",
		"100",
		"Number of lines to show from each service",
	)

	logsCmd.Flags().BoolVar(
		&logsTimestamps,
		"timestamps",
		false,
		"Show timestamps",
	)

	rootCmd.AddCommand(logsCmd)
}
