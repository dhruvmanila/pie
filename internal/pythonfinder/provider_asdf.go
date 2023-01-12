package pythonfinder

import (
	"os"
	"path/filepath"

	"github.com/dhruvmanila/pyvenv/internal/pathutil"
)

// asdfProvider is a Provider that finds Python executables in the
// asdf installation.
type asdfProvider struct {
	// root is the root directory of the asdf installation.
	root string
}

// newAsdfProvider returns a new asdfProvider.
//
// It will return nil if asdf is not installed. This is deduced by checking
// the environment variable ASDF_DATA_DIR, fallback to the default asdf
// installation directory.
func newAsdfProvider() *asdfProvider {
	root := os.Getenv("ASDF_DATA_DIR")
	if root == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil
		}
		root = filepath.Join(homeDir, ".asdf")
	}
	if !pathutil.IsDir(root) {
		return nil
	}
	return &asdfProvider{root: root}
}

func (p *asdfProvider) Executables() ([]string, error) {
	pythonDir := filepath.Join(p.root, "installs", "python")
	if !pathutil.IsDir(pythonDir) {
		return nil, nil
	}

	entries, err := os.ReadDir(pythonDir)
	if err != nil {
		return nil, err
	}

	var executables []string

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		binDir := filepath.Join(pythonDir, entry.Name(), "bin")
		execs, err := execsInPath(binDir)
		if err != nil {
			return nil, err
		}
		executables = append(executables, execs...)
	}

	return executables, nil
}
