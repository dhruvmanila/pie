package pythonfinder

import (
	"fmt"
	"os/exec"
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
