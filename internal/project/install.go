package project

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func Install(projectDir, projectType string) error {
	assets, err := AssetsFor(projectType)
	if err != nil {
		return err
	}

	dir := filepath.Join(projectDir, ".workplace", "nginx")

	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create nginx directory: %w", err)
	}

	path := filepath.Join(dir, "default.conf")

	file, err := os.OpenFile(
		path,
		os.O_WRONLY|os.O_CREATE|os.O_EXCL,
		0644,
	)
	if err != nil {
		if os.IsExist(err) {
			return nil
		}

		return fmt.Errorf("create nginx configuration: %w", err)
	}

	if _, err := file.Write(assets.NginxConfig); err != nil {
		closeErr := file.Close()
		removeErr := os.Remove(path)

		return errors.Join(
			fmt.Errorf("write nginx configuration: %w", err),
			closeErr,
			removeErr,
		)
	}

	if err := file.Close(); err != nil {
		removeErr := os.Remove(path)

		return errors.Join(
			fmt.Errorf("close nginx configuration: %w", err),
			removeErr,
		)
	}

	return nil
}
