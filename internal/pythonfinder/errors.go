package pythonfinder

import "errors"

var (
	// ErrInvalidVersion is returned when a version string is not valid.
	// This could come from the version string provided by the user, or
	// the one returned by the Python executable. Currently, the only valid
	// version string is of the form "major.minor.patch".
	ErrInvalidVersion = errors.New("invalid version")

	// ErrVersionNotFound is returned when either the version provided by
	// the user is not found, or there is no version of Python installed
	// on the system.
	ErrVersionNotFound = errors.New("version does not exist")
)
