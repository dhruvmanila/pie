package pythonfinder

import (
	"regexp"
	"strconv"
)

// versionRegex is a regular expression that matches Python version strings.
//
// This only matches the final version numbers which are of the form
// "major.minor.patch".
var versionRegex = regexp.MustCompile(`^(\d+)\.(\d+)\.(\d+)$`)

// parseVersion parses a Python version string of the form "major.minor.patch".
//
// If the version string is not valid i.e. it does not match the expected
// format, then ErrInvalidVersion is returned. Other error types include
// the errors returned by [strconv.ParseInt].
func parseVersion(version string) (*VersionInfo, error) {
	matches := versionRegex.FindStringSubmatch(version)
	if matches == nil {
		return nil, ErrInvalidVersion
	}
	major, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil, err
	}
	minor, err := strconv.Atoi(matches[2])
	if err != nil {
		return nil, err
	}
	patch, err := strconv.Atoi(matches[3])
	if err != nil {
		return nil, err
	}
	return &VersionInfo{
		Major: major,
		Minor: minor,
		Patch: patch,
	}, nil
}
