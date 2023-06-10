package pythonfinder

import (
	"os/exec"
	"runtime"

	pep440Version "github.com/aquasecurity/go-pep440-version"
)

// finderStrategy is the strategy used by the finder to find Python executables.
//
// This is decided as per the given version.
type finderStrategy int

const (
	// findAll finds all Python executables.
	findAll finderStrategy = iota

	// findFirst finds the first Python executable.
	findFirst

	// findExact finds the Python executable which matches the given version
	// exactly.
	findExact

	// findGlob finds the Python executable which matches the given version
	// using glob matching.
	findGlob
)

func (s finderStrategy) String() string {
	switch s {
	case findAll:
		return "findAll"
	case findFirst:
		return "findFirst"
	case findExact:
		return "findExact"
	case findGlob:
		return "findMax"
	default:
		return "unknown"
	}
}

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
// The strategy used to find the Python version is decided as per the given
// version using the following rules:
//  1. If the given version is empty, find the first Python version.
//  2. If the given version is a final release and not a complete version, find
//     the max version. For example, if the given version is 3.11, find the max
//     version among all the Python versions which match 3.11.*.
//  3. Otherwise, find the Python version which matches the given version exactly.
//
// A final release is a version which is not a pre-release, post-release, or
// developmental release.
//
// A complete version is a version which has all the version components and
// is a final release. For example, 3.11.2 is a complete version but 3.11 is
// not.
func (f *finder) Find(version string) (*PythonExecutable, error) {
	var strategy finderStrategy
	var versionInfo *pep440Version.Version

	if version == "" {
		strategy = findFirst
	} else {
		v, err := pep440Version.Parse(version)
		if err != nil {
			return nil, err
		}
		versionInfo = &v
		if isFinalRelease(versionInfo) && !isCompleteVersion(versionInfo) {
			strategy = findGlob
		} else {
			strategy = findExact
		}
	}

	versions, err := f.find(versionInfo, strategy)
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
	return f.find(nil, findAll)
}

func (f *finder) find(versionInfo *pep440Version.Version, strategy finderStrategy) ([]*PythonExecutable, error) {
	var versions []*PythonExecutable
	var maxVersion *PythonExecutable

	var specifier *pep440Version.Specifiers
	if strategy == findGlob {
		// If strategy is findMax, versionInfo is guaranteed to be non-nil and
		// not a complete version.
		s, err := pep440Version.NewSpecifiers("== " + getGlobVersion(versionInfo))
		if err != nil {
			return nil, err
		}
		specifier = &s
	}

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
				switch err.(type) {
				case *exec.Error, *exec.ExitError:
					// The file could not be classified as an executable or
					// the execution failed. In both cases, we just ignore
					// the error and continue.
					continue
				}
				return nil, err
			}

			switch strategy {
			case findFirst:
				versions = append(versions, pythonExecutable)
				break ProviderLoop
			case findAll:
				versions = append(versions, pythonExecutable)
			default:
				ordering := pythonExecutable.Version.Compare(*versionInfo)
				switch strategy {
				case findExact:
					if ordering == 0 {
						versions = append(versions, pythonExecutable)
						break ProviderLoop
					}
				case findGlob:
					if ordering < 0 || !isFinalRelease(pythonExecutable.Version) {
						continue
					}
					if specifier.Check(*pythonExecutable.Version) {
						if maxVersion == nil {
							maxVersion = pythonExecutable
						} else if pythonExecutable.Version.GreaterThan(*maxVersion.Version) {
							maxVersion = pythonExecutable
						}
					}
				}
			}
		}
	}

	if maxVersion != nil {
		versions = append(versions, maxVersion)
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
