package project

import (
	"embed"
	"fmt"
	"io/fs"
	"strings"
)

//go:embed templates
var templates embed.FS

type Assets struct {
	NginxConfig []byte
}

func AssetsFor(projectType string) (Assets, error) {
	projectType = strings.ToLower(strings.TrimSpace(projectType))

	if projectType == "" {
		return Assets{}, fmt.Errorf("project type is empty")
	}

	path := fmt.Sprintf("templates/%s/nginx.conf", projectType)

	content, err := fs.ReadFile(templates, path)
	if err != nil {
		if !strings.Contains(err.Error(), "file does not exist") {
			return Assets{}, fmt.Errorf(
				"read nginx template for %q: %w",
				projectType,
				err,
			)
		}

		return Assets{}, fmt.Errorf(
			"unsupported project type %q: nginx template not found",
			projectType,
		)
	}

	return Assets{
		NginxConfig: content,
	}, nil
}
