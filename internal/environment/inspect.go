package environment

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mygento/workplace/internal/config"
)

// Snapshot contains all information collected about the host and project.
type Snapshot struct {
	CollectedAt string          `json:"collected_at"`
	ProjectDir  string          `json:"project_dir"`
	Host        HostInfo        `json:"host"`
	PHP         PHPInfo         `json:"php"`
	Composer    ToolInfo        `json:"composer"`
	Node        NodeInfo        `json:"node"`
	Docker      DockerInfo      `json:"docker"`
	Containers  []ContainerInfo `json:"containers"`
	Project     ProjectInfo     `json:"project"`
	Config      ConfigInfo      `json:"config"`
	Findings    []Finding       `json:"findings"`
}

// HostInfo describes the operating system.
type HostInfo struct {
	OS           string `json:"os"`
	Architecture string `json:"architecture"`
	Version      string `json:"version,omitempty"`
	Shell        string `json:"shell,omitempty"`
}

// PHPInfo describes the default PHP executable and discovered PHP versions.
type PHPInfo struct {
	Installed bool     `json:"installed"`
	Path      string   `json:"path,omitempty"`
	Version   string   `json:"version,omitempty"`
	Ini       string   `json:"ini,omitempty"`
	Available []string `json:"available"`
	Error     string   `json:"error,omitempty"`
}

// ToolInfo describes an executable and its version.
type ToolInfo struct {
	Installed bool   `json:"installed"`
	Path      string `json:"path,omitempty"`
	Version   string `json:"version,omitempty"`
	Error     string `json:"error,omitempty"`
}

// NodeInfo describes Node.js and npm.
type NodeInfo struct {
	Node ToolInfo `json:"node"`
	NPM  ToolInfo `json:"npm"`
}

// DockerInfo describes Docker CLI, daemon, and Compose availability.
type DockerInfo struct {
	ClientInstalled  bool     `json:"client_installed"`
	ClientPath       string   `json:"client_path,omitempty"`
	ClientVersion    string   `json:"client_version,omitempty"`
	DaemonAvailable  bool     `json:"daemon_available"`
	DaemonVersion    string   `json:"daemon_version,omitempty"`
	DaemonError      string   `json:"daemon_error,omitempty"`
	ComposeInstalled bool     `json:"compose_installed"`
	ComposeVersion   string   `json:"compose_version,omitempty"`
	ComposeCommand   []string `json:"compose_command,omitempty"`
	ComposeError     string   `json:"compose_error,omitempty"`
}

// ContainerInfo describes a running Docker container.
type ContainerInfo struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Image  string `json:"image"`
	Status string `json:"status"`
}

// ProjectInfo describes the current application and its manifests.
type ProjectInfo struct {
	Root            string            `json:"root"`
	Name            string            `json:"name,omitempty"`
	Type            []string          `json:"type,omitempty"`
	ComposerJSON    ManifestInfo      `json:"composer_json"`
	PackageJSON     ManifestInfo      `json:"package_json"`
	ComposerLock    bool              `json:"composer_lock"`
	PackageLock     bool              `json:"package_lock"`
	YarnLock        bool              `json:"yarn_lock"`
	PnpmLock        bool              `json:"pnpm_lock"`
	ComposerName    string            `json:"composer_name,omitempty"`
	NPMName         string            `json:"npm_name,omitempty"`
	Dependencies    []string          `json:"dependencies,omitempty"`
	DevDependencies []string          `json:"dev_dependencies,omitempty"`
	NPMScripts      map[string]string `json:"npm_scripts,omitempty"`
}

// ManifestInfo describes the presence and validity of a manifest.
type ManifestInfo struct {
	Exists bool   `json:"exists"`
	Valid  bool   `json:"valid"`
	Error  string `json:"error,omitempty"`
}

// ConfigInfo describes the configuration files used and resolved config.
type ConfigInfo struct {
	Files    []string       `json:"files"`
	Resolved *config.Config `json:"resolved"`
}

// Finding represents an environment diagnostic or warning.
type Finding struct {
	Level   string `json:"level"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

// commandResult stores the output of a system command.
type commandResult struct {
	Output string
	Err    error
}

// Collect inspects the local host and project without modifying either.
//
// The project directory defaults to the current working directory.
// Individual tools may be missing or unavailable; these conditions are
// recorded in the snapshot rather than causing collection to fail.
//
// Collect returns an error when the project directory cannot be resolved
// or is not an accessible directory.
func Collect(ctx context.Context, projectDir string) (*Snapshot, error) {
	if strings.TrimSpace(projectDir) == "" {
		projectDir = "."
	}

	root, err := filepath.Abs(projectDir)
	if err != nil {
		return nil, fmt.Errorf("resolve project directory: %w", err)
	}

	info, err := os.Stat(root)
	if err != nil {
		return nil, fmt.Errorf("inspect project directory: %w", err)
	}

	if !info.IsDir() {
		return nil, fmt.Errorf("project path is not a directory: %s", root)
	}

	snapshot := &Snapshot{
		CollectedAt: time.Now().Format(time.RFC3339),
		ProjectDir:  root,
		Containers:  []ContainerInfo{},
		Findings:    []Finding{},
	}

	// Host and installed development tools.
	snapshot.Host = collectHost()
	snapshot.PHP = collectPHP(ctx, root)
	snapshot.Composer = collectComposer(ctx, root)
	snapshot.Node = collectNode(ctx, root)

	// Docker CLI, daemon, Compose, and running containers.
	snapshot.Docker = collectDocker(
		ctx,
		root,
		&snapshot.Containers,
	)

	// Project manifests and lock files.
	snapshot.Project = collectProject(root)

	// Resolve configuration using project-type defaults and local overrides.
	resolvedConfig, configErr := config.NewLoader(root).Load()
	if configErr != nil {
		snapshot.Config = ConfigInfo{
			Files: configFiles(root),
		}

		snapshot.Findings = append(snapshot.Findings, Finding{
			Level:   "warning",
			Code:    "config_load_failed",
			Message: configErr.Error(),
		})
	} else {
		snapshot.Config = ConfigInfo{
			Files:    configFiles(root),
			Resolved: &resolvedConfig,
		}

		// The resolved config is the source of truth for the project type.
		snapshot.Project.Type = []string{resolvedConfig.Type}
	}

	// Resolve the project name using the existing precedence.
	name, err := ResolveProjectName(root)
	if err != nil {
		snapshot.Project.Name = filepath.Base(filepath.Clean(root))
	} else {
		snapshot.Project.Name = name
	}

	// Produce actionable diagnostics from the collected information.
	collectFindings(snapshot)

	return snapshot, nil
}

// configFiles returns the paths of configuration files present in the project,
// ordered from shared project configuration to local overrides.
func configFiles(root string) []string {
	candidates := []string{
		filepath.Join(root, ".workplace.yaml"),
		filepath.Join(root, ".workplace", "config.local.yaml"),
	}

	files := make([]string, 0, len(candidates))

	for _, path := range candidates {
		info, err := os.Stat(path)
		if err != nil || !info.Mode().IsRegular() {
			continue
		}

		files = append(files, path)
	}

	return files
}
