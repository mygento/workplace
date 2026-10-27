package environment

func collectFindings(s *Snapshot) {
	if !s.PHP.Installed && len(s.PHP.Available) == 0 {
		s.Findings = append(s.Findings, Finding{
			Level:   "warning",
			Code:    "php_missing",
			Message: "No working PHP executable was found in PATH.",
		})
	}

	if !s.Composer.Installed {
		s.Findings = append(s.Findings, Finding{
			Level:   "warning",
			Code:    "composer_missing",
			Message: "Composer was not found or could not be executed.",
		})
	}

	if !s.Docker.ClientInstalled {
		s.Findings = append(s.Findings, Finding{
			Level:   "info",
			Code:    "docker_missing",
			Message: "Docker CLI was not found in PATH.",
		})
	} else if !s.Docker.DaemonAvailable {
		s.Findings = append(s.Findings, Finding{
			Level: "warning",
			Code:  "docker_daemon_unavailable",
			Message: "Docker CLI is installed, but the daemon is unavailable: " +
				s.Docker.DaemonError,
		})
	}

	if !s.Docker.ComposeInstalled {
		s.Findings = append(s.Findings, Finding{
			Level:   "info",
			Code:    "compose_missing",
			Message: "Docker Compose was not detected.",
		})
	}
}
