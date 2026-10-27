package environment

import (
	"context"
	"fmt"
	"os/exec"
	"sort"
	"strings"
)

func collectPHP(ctx context.Context, dir string) PHPInfo {
	var result PHPInfo
	names := []string{"php"}

	// Includes php5.6, php7.0-php7.9 and php8.0-php8.9.
	for major := 5; major <= 8; major++ {
		startMinor := 0
		maxMinor := 9

		if major == 5 {
			startMinor, maxMinor = 6, 6
		}

		for minor := startMinor; minor <= maxMinor; minor++ {
			names = append(names, fmt.Sprintf("php%d.%d", major, minor))
		}
	}

	seen := make(map[string]bool)

	for _, name := range names {
		path, err := exec.LookPath(name)
		if err != nil || seen[path] {
			continue
		}
		seen[path] = true

		version := run(ctx, dir, path, "-r", "echo PHP_VERSION;")
		if version.Err != nil || version.Output == "" {
			continue
		}

		result.Available = append(
			result.Available,
			fmt.Sprintf("%s (%s)", path, version.Output),
		)

		if name == "php" {
			result.Installed = true
			result.Path = path
			result.Version = version.Output

			ini := run(ctx, dir, path, "--ini")
			if ini.Err == nil {
				for _, line := range strings.Split(ini.Output, "\n") {
					if strings.HasPrefix(line, "Loaded Configuration File:") {
						result.Ini = strings.TrimSpace(
							strings.TrimPrefix(
								line, "Loaded Configuration File:",
							),
						)
						break
					}
				}
			}
		}
	}

	sort.Strings(result.Available)

	if !result.Installed {
		result.Error = "default php executable not found; versioned executables may still be available"
	}

	return result
}
