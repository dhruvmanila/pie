package python

import (
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

// VersionNotFoundError is returned by VersionLookup when it fails
// to find an executable for the provided Python version.
type VersionNotFoundError struct {
	// Version is the Python version for which the error occurred.
	Version string
}

func (e *VersionNotFoundError) Error() string {
	return e.Version + ": version does not exist"
}

// constructCommandName is used to contruct the Python command name
// from the given version string.
//
// The output is OS specific where windows commands are of the form
// 'pythonXY.exe' and unix commands are of the form 'pythonX.Y' where
// X and Y are the major and minor part of the version.
func constructCommandName(version string) string {
	if version == "" {
		return DefaultExec
	}
	versionParts := strings.Split(version, ".")
	return fmt.Sprintf(ExecVersionFormat, versionParts[0], versionParts[1])
}

// VersionLookup searches for the Python executable for the given version
// and returns the absolute path to the file.
//
// If no version is provided, then the global Python executable will be used.
// The lookup is based on the major and minor parts of the version and if
// there are multiple executables available on PATH, this will only check
// against the first one.
//
// If it's unable to find an exectuable that matches the given version, the
// error is of type *VersionNotFoundError. Other error types may be returned
// for other situations.
func VersionLookup(version string) (string, error) {
	cmd := exec.Command(constructCommandName(version), "--version")

	output, err := cmd.Output()
	if err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			return "", &VersionNotFoundError{version}
		}
		return "", err
	}

	realpath, err := filepath.EvalSymlinks(cmd.Path)
	if err != nil {
		return "", err
	}
	if version == "" {
		return realpath, nil
	}

	// Output: "Python <version><LF/CRLF>"
	actualVersion := strings.Split(string(output), " ")[1]
	actualVersion = strings.TrimRightFunc(actualVersion, func(r rune) bool {
		return r == '\n' || r == '\r'
	})
	if version != actualVersion {
		return "", &VersionNotFoundError{version}
	}

	return realpath, nil
}
