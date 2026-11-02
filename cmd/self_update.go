package cmd

import (
	"fmt"

	"github.com/creativeprojects/go-selfupdate"
	"github.com/spf13/cobra"
)

const (
	selfUpdateOwner = "mygento"
	selfUpdateRepo  = "workplace"
)

var selfUpdateCmd = &cobra.Command{
	Use:     "self-update",
	Aliases: []string{"update"},
	Short:   "Update Workplace to the latest version",
	Long:    "Check GitHub for the latest Workplace release and update the current executable.",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		fmt.Printf("Current version: v%s\n", version)
		fmt.Println("Checking for updates...")

		updater, err := selfupdate.NewUpdater(selfupdate.Config{
			Validator: &selfupdate.ChecksumValidator{
				UniqueFilename: "checksums.txt",
			},
		})
		if err != nil {
			return fmt.Errorf("initialize updater: %w", err)
		}

		repository := selfupdate.NewRepositorySlug(
			selfUpdateOwner,
			selfUpdateRepo,
		)

		release, found, err := updater.DetectLatest(ctx, repository)
		if err != nil {
			return fmt.Errorf("check latest release: %w", err)
		}

		if !found {
			return fmt.Errorf(
				"no compatible release found for %s/%s",
				"current operating system",
				"current architecture",
			)
		}

		if release.LessOrEqual(version) {
			fmt.Printf("Workplace is already up to date (v%s)\n", version)
			return nil
		}

		fmt.Printf("New version available: v%s\n", release.Version())

		executable, err := selfupdate.ExecutablePath()
		if err != nil {
			return fmt.Errorf("locate Workplace executable: %w", err)
		}

		fmt.Println("Downloading and installing update...")

		if err := updater.UpdateTo(ctx, release, executable); err != nil {
			return fmt.Errorf("update Workplace: %w", err)
		}

		fmt.Printf(
			"Successfully updated Workplace: v%s -> v%s\n",
			version,
			release.Version(),
		)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(selfUpdateCmd)
}
