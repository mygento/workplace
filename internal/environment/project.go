package environment

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type composerManifest struct {
	Name       string            `json:"name"`
	Require    map[string]string `json:"require"`
	RequireDev map[string]string `json:"require-dev"`
}

type npmManifest struct {
	Name            string            `json:"name"`
	Scripts         map[string]string `json:"scripts"`
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
	Workplace       json.RawMessage   `json:"workplace"`
}

func collectProject(root string) ProjectInfo {
	p := ProjectInfo{Root: root}

	composerPath := filepath.Join(root, "composer.json")
	data, err := os.ReadFile(composerPath)
	p.ComposerJSON.Exists = err == nil

	if err == nil {
		var manifest composerManifest
		if jsonErr := json.Unmarshal(data, &manifest); jsonErr != nil {
			p.ComposerJSON.Error = jsonErr.Error()
		} else {
			p.ComposerJSON.Valid = true
			p.ComposerName = manifest.Name

			for name := range manifest.Require {
				if name != "php" && name != "composer-plugin-api" {
					p.Dependencies = append(p.Dependencies, name)
				}
			}
			for name := range manifest.RequireDev {
				p.DevDependencies = append(p.DevDependencies, name)
			}

			if _, ok := manifest.Require["magento/product-community-edition"]; ok {
				p.Type = appendUnique(p.Type, "magento2")
			}
			if _, ok := manifest.Require["magento/product-enterprise-edition"]; ok {
				p.Type = appendUnique(p.Type, "magento2")
			}
			if _, ok := manifest.Require["magento/magento2ce"]; ok {
				p.Type = appendUnique(p.Type, "magento2")
			}
			if _, ok := manifest.Require["symfony/framework-bundle"]; ok {
				p.Type = appendUnique(p.Type, "symfony")
			}
		}
	} else if !os.IsNotExist(err) {
		p.ComposerJSON.Error = err.Error()
	}

	packagePath := filepath.Join(root, "package.json")
	data, err = os.ReadFile(packagePath)
	p.PackageJSON.Exists = err == nil

	if err == nil {
		var manifest npmManifest
		if jsonErr := json.Unmarshal(data, &manifest); jsonErr != nil {
			p.PackageJSON.Error = jsonErr.Error()
		} else {
			p.PackageJSON.Valid = true
			p.NPMName = manifest.Name
			p.NPMScripts = manifest.Scripts
		}
	} else if !os.IsNotExist(err) {
		p.PackageJSON.Error = err.Error()
	}

	p.ComposerLock = fileExists(filepath.Join(root, "composer.lock"))
	p.PackageLock = fileExists(filepath.Join(root, "package-lock.json"))
	p.YarnLock = fileExists(filepath.Join(root, "yarn.lock"))
	p.PnpmLock = fileExists(filepath.Join(root, "pnpm-lock.yaml"))

	for _, dependency := range p.Dependencies {
		if strings.HasPrefix(dependency, "magento/") {
			p.Type = appendUnique(p.Type, "magento2")
		}
		if strings.HasPrefix(dependency, "symfony/") {
			p.Type = appendUnique(p.Type, "symfony")
		}
	}

	sort.Strings(p.Type)
	sort.Strings(p.Dependencies)
	sort.Strings(p.DevDependencies)

	return p
}
