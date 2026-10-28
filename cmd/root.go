package cmd

import (
	"fmt"
	"os"

	workplaceruntime "github.com/mygento/workplace/internal/runtime"
	"github.com/spf13/cobra"
)

var (
	verbose bool

	rootCmd = &cobra.Command{
		Use:   "workplace",
		Short: "workplace is a cli tool for devs",
		Long:  "workplace is a cli tool for devs",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			workplaceruntime.SetVerbose(verbose)
		},
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error executing workplace: %s\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(
		&verbose,
		"verbose",
		"v",
		false,
		"show commands executed by workplace",
	)
}
