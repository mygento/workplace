package cmd

import (
	"github.com/mygento/workplace/internal/server"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(startCmd)
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a workplace environment",
	Long:  `Start a workplace environment`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return server.Start()
	},
}
