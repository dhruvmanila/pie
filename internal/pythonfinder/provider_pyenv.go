package pythonfinder

import (
	"os"
	"path/filepath"
)

// pyenvProvider is a Provider that finds Python executables in the
// pyenv installation.
type pyenvProvider struct {
	// root is the root directory of the pyenv installation.
	root string
}

// newPyenvProvider returns a new pyenvProvider.
//
// It will return nil if pyenv is not installed. This is deduced by checking
// the environment variable PYENV_ROOT, fallback to the default pyenv
// installation directory.
func newPyenvProvider() *pyenvProvider {
	root := os.Getenv("PYENV_ROOT")
	if root == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil
		}
		root = filepath.Join(homeDir, ".pyenv")
	}
	if stat, err := os.Stat(root); err != nil || !stat.IsDir() {
		return nil
	}
	return &pyenvProvider{root: root}
}

func (p *pyenvProvider) Executables() ([]string, error) {
	versionDir := filepath.Join(p.root, "versions")
	if !isDir(versionDir) {
		return nil, nil
	}

	entries, err := os.ReadDir(versionDir)
	if err != nil {
		return nil, err
	}

	var executables []string

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		binDir := filepath.Join(versionDir, entry.Name(), "bin")
		if !isDir(binDir) {
			continue
		}
		execs, err := execsInPath(binDir)
		if err != nil {
			return nil, err
		}
		executables = append(executables, execs...)
	}

	return executables, nil
}
