package cmd

import (
	"errors"
	"fmt"

	"github.com/mygento/workplace/internal/environment"
	"github.com/spf13/cobra"
)

var migrateConfigCmd = &cobra.Command{
	Use:   "migrate-config",
	Short: "Migrate legacy package.json workplace configuration to YAML",
	Long:  `Extract the workplace configuration from package.json and write it to .workplace.yaml.`,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, err := cmd.Flags().GetString("dir")
		if err != nil {
			return err
		}

		overwrite, err := cmd.Flags().GetBool("overwrite")
		if err != nil {
			return err
		}

		result, err := environment.MigrateConfig(dir, overwrite)
		if err != nil {
			if errors.Is(err, environment.ErrTargetConfigExists) {
				return fmt.Errorf(
					".workplace.yaml already exists; inspect it first or use --overwrite",
				)
			}

			if errors.Is(err, environment.ErrLegacyConfigNotFound) {
				return fmt.Errorf(
					"no legacy workplace configuration found in package.json",
				)
			}

			return err
		}

		out := cmd.OutOrStdout()

		if _, err := fmt.Fprintln(out, "Configuration migrated successfully."); err != nil {
			return err
		}
		if _, err := fmt.Fprintf(out, "Source:       %s\n", result.Source); err != nil {
			return err
		}
		if _, err := fmt.Fprintf(out, "Destination:  %s\n", result.Destination); err != nil {
			return err
		}
		if _, err := fmt.Fprintf(out, "Project type: %s\n", result.Type); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	migrateConfigCmd.Flags().StringP(
		"dir", "d", ".", "Project directory",
	)
	migrateConfigCmd.Flags().Bool(
		"overwrite", false, "Overwrite an existing .workplace.yaml",
	)

	rootCmd.AddCommand(migrateConfigCmd)
}
