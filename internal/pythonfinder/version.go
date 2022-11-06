package pythonfinder

import (
	"fmt"
	"os/exec"
	"strings"
)

// VersionInfo contains information about a Python version. It provides
// information about the major, minor, and patch version numbers.
type VersionInfo struct {
	Major int
	Minor int
	Patch int
}

// Matches returns true if the other version matches the current version.
func (v *VersionInfo) Matches(other *VersionInfo) bool {
	return v.Major == other.Major && v.Minor == other.Minor && v.Patch == other.Patch
}

func (v *VersionInfo) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

// PythonVersion contains information about a Python version.
//
// This is similar to VersionInfo, but also contains the path to the Python
// executable.
type PythonVersion struct {
	*VersionInfo

	// Executable is the absolute path to the Python executable.
	Executable string
}

// newPythonVersion creates a new PythonVersion from the given executable.
//
// This function returns ErrInvalidVersion if the executable Python version
// is not supported by this tool. It also returns errors from the exec package.
func newPythonVersion(executable string) (*PythonVersion, error) {
	versionInfo, err := getVersionInfo(executable)
	if err != nil {
		return nil, err
	}
	return &PythonVersion{VersionInfo: versionInfo, Executable: executable}, nil
}

func (v *PythonVersion) String() string {
	return fmt.Sprintf("%s (%s)", v.VersionInfo, v.Executable)
}

// getVersionInfo returns the version information for the given Python
// executable.
func getVersionInfo(executable string) (*VersionInfo, error) {
	cmd := exec.Command(executable, "--version")

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	// Output: "Python <version><LF/CRLF>"
	_, version, found := strings.Cut(string(output), " ")
	if !found {
		return nil, fmt.Errorf("Unable to parse Python version: %q", output)
	}
	version = strings.TrimRight(version, "\r\n")

	versionInfo, err := parseVersion(version)
	if err != nil {
		return nil, err
	}

	return versionInfo, nil
}
