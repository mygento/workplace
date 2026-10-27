package cmd

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/mygento/workplace/internal/environment"
	"github.com/spf13/cobra"
)

var inspectCmd = &cobra.Command{
	Use:   "inspect",
	Short: "Inspect the local development environment",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, err := cmd.Flags().GetString("dir")
		if err != nil {
			return err
		}

		asJSON, err := cmd.Flags().GetBool("json")
		if err != nil {
			return err
		}

		s, err := environment.Collect(cmd.Context(), dir)
		if err != nil {
			return err
		}

		if asJSON {
			encoder := json.NewEncoder(cmd.OutOrStdout())
			encoder.SetIndent("", "  ")
			return encoder.Encode(s)
		}

		return printInspection(cmd, s)
	},
}

func init() {
	inspectCmd.Flags().StringP("dir", "d", ".", "Project directory")
	inspectCmd.Flags().Bool("json", false, "Output JSON")
	rootCmd.AddCommand(inspectCmd)
}

type inspectionWriter interface {
	Write([]byte) (int, error)
}

func printInspection(cmd *cobra.Command, s *environment.Snapshot) error {
	out := cmd.OutOrStdout()

	if err := writeLine(out, "Workplace environment inspection"); err != nil {
		return err
	}
	if err := writeLine(out, strings.Repeat("=", 34)); err != nil {
		return err
	}
	if err := writef(out, "Collected: %s\n", s.CollectedAt); err != nil {
		return err
	}
	if err := writef(out, "Directory: %s\n", s.ProjectDir); err != nil {
		return err
	}

	if err := writeLine(out, "\n[Host]"); err != nil {
		return err
	}
	if err := writef(out, "OS: %s\n", s.Host.OS); err != nil {
		return err
	}
	if err := writef(out, "Architecture: %s\n", s.Host.Architecture); err != nil {
		return err
	}
	if err := printValue(out, "OS version", s.Host.Version); err != nil {
		return err
	}
	if err := printValue(out, "Shell", s.Host.Shell); err != nil {
		return err
	}

	if err := writeLine(out, "\n[PHP]"); err != nil {
		return err
	}
	if err := printBool(out, "Default PHP installed", s.PHP.Installed); err != nil {
		return err
	}
	if err := printValue(out, "Executable", s.PHP.Path); err != nil {
		return err
	}
	if err := printValue(out, "Version", s.PHP.Version); err != nil {
		return err
	}
	if err := printValue(out, "php.ini", s.PHP.Ini); err != nil {
		return err
	}
	for _, version := range s.PHP.Available {
		if err := writef(out, "  - %s\n", version); err != nil {
			return err
		}
	}
	if err := printValue(out, "Error", s.PHP.Error); err != nil {
		return err
	}

	if err := writeLine(out, "\n[Composer]"); err != nil {
		return err
	}
	if err := printBool(out, "Installed", s.Composer.Installed); err != nil {
		return err
	}
	if err := printValue(out, "Executable", s.Composer.Path); err != nil {
		return err
	}
	if err := printValue(out, "Version", s.Composer.Version); err != nil {
		return err
	}
	if err := printValue(out, "Error", s.Composer.Error); err != nil {
		return err
	}

	if err := writeLine(out, "\n[Node.js]"); err != nil {
		return err
	}
	if err := printTool(out, "Node.js", s.Node.Node); err != nil {
		return err
	}
	if err := printTool(out, "npm", s.Node.NPM); err != nil {
		return err
	}

	if err := writeLine(out, "\n[Docker]"); err != nil {
		return err
	}
	if err := printBool(out, "CLI installed", s.Docker.ClientInstalled); err != nil {
		return err
	}
	if err := printValue(out, "Executable", s.Docker.ClientPath); err != nil {
		return err
	}
	if err := printValue(out, "Client version", s.Docker.ClientVersion); err != nil {
		return err
	}
	if err := printBool(out, "Daemon available", s.Docker.DaemonAvailable); err != nil {
		return err
	}
	if err := printValue(out, "Daemon version", s.Docker.DaemonVersion); err != nil {
		return err
	}
	if err := printValue(out, "Daemon error", s.Docker.DaemonError); err != nil {
		return err
	}
	if err := printBool(out, "Compose installed", s.Docker.ComposeInstalled); err != nil {
		return err
	}
	if err := printValue(out, "Compose version", s.Docker.ComposeVersion); err != nil {
		return err
	}
	if err := printValue(out, "Compose command", strings.Join(s.Docker.ComposeCommand, " ")); err != nil {
		return err
	}
	if err := printValue(out, "Compose error", s.Docker.ComposeError); err != nil {
		return err
	}

	if err := writef(out, "\n[Running containers: %d]\n", len(s.Containers)); err != nil {
		return err
	}
	for _, c := range s.Containers {
		if err := writef(out, "- %s | %s | %s | %s\n", c.Name, c.Image, c.Status, c.ID); err != nil {
			return err
		}
	}

	if err := writeLine(out, "\n[Project]"); err != nil {
		return err
	}
	if err := printValue(out, "Name", s.Project.Name); err != nil {
		return err
	}
	if err := printValue(out, "Root", s.Project.Root); err != nil {
		return err
	}
	if err := writef(out, "Type: %s\n", strings.Join(s.Project.Type, ", ")); err != nil {
		return err
	}
	if err := printManifest(out, "composer.json", s.Project.ComposerJSON); err != nil {
		return err
	}
	if err := printManifest(out, "package.json", s.Project.PackageJSON); err != nil {
		return err
	}
	if err := printBool(out, "composer.lock", s.Project.ComposerLock); err != nil {
		return err
	}
	if err := printBool(out, "package-lock.json", s.Project.PackageLock); err != nil {
		return err
	}
	if err := printBool(out, "yarn.lock", s.Project.YarnLock); err != nil {
		return err
	}
	if err := printBool(out, "pnpm-lock.yaml", s.Project.PnpmLock); err != nil {
		return err
	}
	if err := printValue(out, "Composer package", s.Project.ComposerName); err != nil {
		return err
	}
	if err := printValue(out, "NPM package", s.Project.NPMName); err != nil {
		return err
	}
	if len(s.Project.Dependencies) > 0 {
		if err := writef(out, "Composer dependencies (%d): %s\n",
			len(s.Project.Dependencies),
			strings.Join(s.Project.Dependencies, ", ")); err != nil {
			return err
		}
	}
	if len(s.Project.DevDependencies) > 0 {
		if err := writef(out, "Composer dev dependencies (%d): %s\n",
			len(s.Project.DevDependencies),
			strings.Join(s.Project.DevDependencies, ", ")); err != nil {
			return err
		}
	}

	if err := printResolvedConfig(out, s); err != nil {
		return err
	}

	if err := writef(out, "\n[Findings: %d]\n", len(s.Findings)); err != nil {
		return err
	}
	for _, finding := range s.Findings {
		if err := writef(out, "- [%s] %s: %s\n",
			strings.ToUpper(finding.Level),
			finding.Code,
			finding.Message); err != nil {
			return err
		}
	}

	return nil
}

func printResolvedConfig(out inspectionWriter, s *environment.Snapshot) error {
	if err := writeLine(out, "\n[Workplace configuration]"); err != nil {
		return err
	}

	if len(s.Config.Files) == 0 {
		if err := writeLine(out, "Files: none"); err != nil {
			return err
		}
	} else {
		if err := writeLine(out, "Files:"); err != nil {
			return err
		}
		for _, path := range s.Config.Files {
			if err := writef(out, "  - %s\n", path); err != nil {
				return err
			}
		}
	}

	if s.Config.Resolved == nil {
		return writeLine(out, "Resolved configuration: unavailable")
	}

	cfg := s.Config.Resolved

	if err := printValue(out, "Name", cfg.Name); err != nil {
		return err
	}
	if err := printValue(out, "Project type", cfg.Type); err != nil {
		return err
	}

	if err := writeLine(out, "Services:"); err != nil {
		return err
	}

	names := make([]string, 0, len(cfg.Services))
	for name := range cfg.Services {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		service := cfg.Services[name]

		if err := writef(out, "  %s:\n", name); err != nil {
			return err
		}
		if err := printValue(out, "    Image", service.Image); err != nil {
			return err
		}
		if service.Port != nil {
			if err := writef(out, "    Port: %d\n", *service.Port); err != nil {
				return err
			}
		}

		if len(service.Env) > 0 {
			if err := writeLine(out, "    Environment:"); err != nil {
				return err
			}

			envNames := make([]string, 0, len(service.Env))
			for key := range service.Env {
				envNames = append(envNames, key)
			}
			sort.Strings(envNames)

			for _, key := range envNames {
				// Do not expose potentially sensitive environment values.
				if err := writef(out, "      %s=<redacted>\n", key); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func writeLine(out inspectionWriter, value string) error {
	_, err := fmt.Fprintln(out, value)
	return err
}

func writef(out inspectionWriter, format string, args ...any) error {
	_, err := fmt.Fprintf(out, format, args...)
	return err
}

func printTool(out inspectionWriter, name string, t environment.ToolInfo) error {
	if err := writef(out, "%s:\n", name); err != nil {
		return err
	}
	if err := printBool(out, "  Installed", t.Installed); err != nil {
		return err
	}
	if err := printValue(out, "  Executable", t.Path); err != nil {
		return err
	}
	if err := printValue(out, "  Version", t.Version); err != nil {
		return err
	}
	return printValue(out, "  Error", t.Error)
}

func printValue(out inspectionWriter, label, value string) error {
	if value == "" {
		return nil
	}
	return writef(out, "%s: %s\n", label, value)
}

func printBool(out inspectionWriter, label string, value bool) error {
	status := "no"
	if value {
		status = "yes"
	}
	return writef(out, "%s: %s\n", label, status)
}

func printManifest(out inspectionWriter, name string, manifest environment.ManifestInfo) error {
	switch {
	case !manifest.Exists:
		return writef(out, "%s: not found\n", name)
	case manifest.Valid:
		return writef(out, "%s: valid JSON\n", name)
	default:
		return writef(out, "%s: invalid JSON (%s)\n", name, manifest.Error)
	}
}
