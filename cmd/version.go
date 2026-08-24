package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of a workplace tool",
	Long:  `All software has versions. This is Workplace's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Workplace v3.0.0")
	},
}
