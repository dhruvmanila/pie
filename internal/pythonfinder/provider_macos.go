package pythonfinder

import (
	"os"
	"path/filepath"

	"github.com/dhruvmanila/pyvenv/internal/pathutil"
)

const versionDir = "/Library/Frameworks/Python.framework/Versions"

// macosProvider is a Provider that finds Python executables in the
// typical macOS installation directory using python.org's installer.
type macosProvider struct{}

// newMacOSProvider returns a new macOSProvider. It returns nil if the
// base installation directory does not exist.
func newMacOSProvider() *macosProvider {
	if !pathutil.IsDir(versionDir) {
		return nil
	}
	return &macosProvider{}
}

func (p *macosProvider) Executables() ([]string, error) {
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
		execs, err := execsInPath(binDir)
		if err != nil {
			return nil, err
		}
		executables = append(executables, execs...)
	}

	return executables, nil
}
