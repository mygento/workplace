package php

import (
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	RuntimeDir string
	SocketPath string
}

func WriteConfig(cfg Config) (string, error) {
	if err := os.MkdirAll(cfg.RuntimeDir, 0755); err != nil {
		return "", fmt.Errorf(
			"create PHP runtime directory: %w",
			err,
		)
	}

	configPath := filepath.Join(
		cfg.RuntimeDir,
		"php-fpm.conf",
	)

	logPath := filepath.Join(
		cfg.RuntimeDir,
		"php-fpm.log",
	)

	pidPath := filepath.Join(
		cfg.RuntimeDir,
		"php-fpm.pid",
	)

	content := fmt.Sprintf(`[global]

daemonize = no
pid = %s
error_log = %s
log_level = notice


[www]

listen = %s
listen.mode = 0666

pm = dynamic
pm.max_children = 10
pm.start_servers = 2
pm.min_spare_servers = 1
pm.max_spare_servers = 3

clear_env = no
catch_workers_output = yes
`,
		pidPath,
		logPath,
		cfg.SocketPath,
	)

	if err := os.WriteFile(
		configPath,
		[]byte(content),
		0644,
	); err != nil {
		return "", fmt.Errorf(
			"write PHP-FPM config: %w",
			err,
		)
	}

	return configPath, nil
}
