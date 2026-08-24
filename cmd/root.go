package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "workplace",
	Short: "workplace is a cli tool for devs",
	Long:  "workplace is a cli tool for devs",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Oops. An error while executing Workplace '%s'\n", err)
		os.Exit(1)
	}
}
