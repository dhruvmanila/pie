package pythonfinder

import (
	"errors"
	"os/exec"
	"runtime"
)

// finder is a Python version finder.
type finder struct {
	providers []Provider
}

// New returns a new Python version finder.
func New() *finder {
	f := &finder{}
	f.setupProviders()
	return f
}

// Find returns the Python version which matches the given version, if provided,
// or the first version found by the providers.
//
// This function returns ErrVersionNotFound if no version matching the given
// version is found or there is no Python version installed on the system.
func (f *finder) Find(version string) (*PythonVersion, error) {
	var versionInfo *VersionInfo

	if version != "" {
		var err error
		versionInfo, err = parseVersion(version)
		if err != nil {
			return nil, err
		}
	}

	for _, p := range f.providers {
		executables, err := p.Executables()
		if err != nil {
			return nil, err
		}

		for _, executable := range executables {
			pythonVersion, err := newPythonVersion(executable)
			if err != nil {
				// The executable could come from a Python version which is not
				// supported by this tool. In this case, we just ignore the error
				// and continue.
				if errors.Is(err, ErrInvalidVersion) {
					continue
				}
				switch err.(type) {
				case *exec.Error, *exec.ExitError:
					// The file could not be classified as an executable or
					// the execution failed. In both cases, we just ignore
					// the error and continue.
					continue
				}
				return nil, err
			}

			if versionInfo != nil {
				if pythonVersion.Matches(versionInfo) {
					return pythonVersion, nil
				}
			} else {
				// If no version was specified by the user, we return the first
				// Python version we find.
				return pythonVersion, nil
			}
		}
	}

	// This either means that the version provided by the user does not exist,
	// or there is no version of Python installed on the system.
	return nil, ErrVersionNotFound
}

func (f *finder) setupProviders() {
	f.providers = append(f.providers, newPathProvider())
	if runtime.GOOS != "windows" {
		if p := newPyenvProvider(); p != nil {
			f.providers = append(f.providers, p)
		}
	}
}
