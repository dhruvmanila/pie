package pythonfinder

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	pep440Version "github.com/aquasecurity/go-pep440-version"
)

var execCommand = exec.Command

// PythonExecutable contains information about a Python executable.
type PythonExecutable struct {
	// Version is the parsed Python version.
	Version *pep440Version.Version

	// Path is the absolute path to the Python executable.
	Path string
}

// newPythonExecutable creates a new PythonExecutable from the given Python
// executable path.
func newPythonExecutable(executable string) (*PythonExecutable, error) {
	versionInfo, err := getPythonVersion(executable)
	if err != nil {
		return nil, err
	}
	return &PythonExecutable{Version: versionInfo, Path: executable}, nil
}

func (v *PythonExecutable) String() string {
	return fmt.Sprintf("%s (%s)", v.Version, v.Path)
}

// getPythonVersion returns the version information for the given Python
// executable.
func getPythonVersion(executable string) (*pep440Version.Version, error) {
	cmd := execCommand(executable, "--version")

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

	versionInfo, err := pep440Version.Parse(version)
	if err != nil {
		return nil, err
	}

	return &versionInfo, nil
}

// isFinalRelease returns true if the given version is a final release, i.e.,
// not a pre-release, post-release or developmental release.
func isFinalRelease(version *pep440Version.Version) bool {
	return !version.IsPreRelease() && !version.IsPostRelease() && version.Local() == ""
}

// versionRegex is a regular expression that matches Python version strings
// of the form "X.Y.Z".
var versionRegex = regexp.MustCompile(`^(\d)\.(\d{1,2})\.(\d{1,2})$`)

// isCompleteVersion returns true if the given version is a complete version,
// i.e., it's of the form X.Y.Z.
//
// This function assumes that the given version is a final release.
func isCompleteVersion(version *pep440Version.Version) bool {
	return versionRegex.MatchString(version.String())
}

// getGlobVersion returns the glob version for the given version if it's not
// a complete version, otherwise it returns the given version.
//
// This function assumes that the given version is a final release.
func getGlobVersion(version *pep440Version.Version) string {
	v := version.String()
	switch len(strings.Split(v, ".")) {
	case 1, 2:
		return fmt.Sprintf("%s.*", v)
	default:
		return v
	}
}
