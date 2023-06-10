package pythonfinder

import "errors"

// ErrVersionNotFound is returned when either the version provided by
// the user is not found, or there is no version of Python installed
// on the system.
var ErrVersionNotFound = errors.New("version does not exist")
