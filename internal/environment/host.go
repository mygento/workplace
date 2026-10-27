package environment

import (
	"os"
	"runtime"
	"strings"
)

func collectHost() HostInfo {
	h := HostInfo{
		OS:           runtime.GOOS,
		Architecture: runtime.GOARCH,
		Shell:        os.Getenv("SHELL"),
	}

	if h.Shell == "" {
		h.Shell = os.Getenv("COMSPEC")
	}

	switch runtime.GOOS {
	case "linux":
		h.Version = readOSRelease()
	case "darwin":
		h.Version = runSimple("sw_vers", "-productVersion")
	case "windows":
		h.Version = runSimple("cmd", "/c", "ver")
	}

	return h
}

func readOSRelease() string {
	data, err := os.ReadFile("/etc/os-release")
	if err != nil {
		return ""
	}

	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "PRETTY_NAME=") {
			return strings.Trim(
				strings.TrimPrefix(line, "PRETTY_NAME="), `"`,
			)
		}
	}
	return ""
}
