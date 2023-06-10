package pythonfinder

import (
	"errors"
	"os/exec"
	"runtime"

	pep440Version "github.com/aquasecurity/go-pep440-version"
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
func (f *finder) Find(version string) (*PythonExecutable, error) {
	versions, err := f.find(version, 1)
	if err != nil {
		return nil, err
	}
	// There is going to be exactly one version in the slice. This is ensured
	// by the find() function. Otherwise, it would return ErrVersionNotFound.
	return versions[0], nil
}

// FindAll returns all the Python versions available on the system which can be
// found by the providers.
func (f *finder) FindAll() ([]*PythonExecutable, error) {
	return f.find("", -1)
}

func (f *finder) find(version string, n int) ([]*PythonExecutable, error) {
	var versionInfo *pep440Version.Version

	if version != "" {
		v, err := pep440Version.Parse(version)
		if err != nil {
			return nil, err
		}
		versionInfo = &v
	}

	var versions []*PythonExecutable

	// seen is a set of Python executables which were already seen by the
	// providers. This is used to avoid returning duplicate Python versions.
	// This contains the absolute path to the Python executable.
	seen := make(map[string]struct{})

ProviderLoop:
	for _, p := range f.providers {
		executables, err := p.Executables()
		if err != nil {
			return nil, err
		}

		for _, executable := range executables {
			if _, ok := seen[executable]; ok {
				continue
			}
			seen[executable] = struct{}{}

			pythonExecutable, err := newPythonExecutable(executable)
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
				if pythonExecutable.Version.Equal(*versionInfo) {
					versions = append(versions, pythonExecutable)
				}
			} else {
				versions = append(versions, pythonExecutable)
			}

			if n > 0 && len(versions) == n {
				break ProviderLoop
			}
		}
	}

	// This either means that the version provided by the user does not exist,
	// or there is no version of Python installed on the system.
	if len(versions) == 0 {
		return nil, ErrVersionNotFound
	}
	return versions, nil
}

func (f *finder) setupProviders() {
	f.providers = append(f.providers, newPathProvider())
	if runtime.GOOS == "darwin" {
		if p := newMacOSProvider(); p != nil {
			f.providers = append(f.providers, p)
		}
	}
	if runtime.GOOS != "windows" {
		if p := newPyenvProvider(); p != nil {
			f.providers = append(f.providers, p)
		}
		if p := newAsdfProvider(); p != nil {
			f.providers = append(f.providers, p)
		}
	}
}
