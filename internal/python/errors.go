package python

// VersionNotFoundError is returned by VersionLookup when it fails
// to find an executable for the provided Python version.
type VersionNotFoundError struct {
	// Version is the Python version for which the error occurred.
	Version string
}

func (e *VersionNotFoundError) Error() string {
	return e.Version + ": version does not exist"
}
