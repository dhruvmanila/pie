package pythonfinder

import (
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// Provider is an interface that provides Python executables.
//
// Each Provider has contextual information that is used to find Python
// executables.
type Provider interface {
	// Executables returns a list of absolute paths to Python executables.
	Executables() ([]string, error)
}

// execsInPath returns a list of Python executables in the given path.
//
// The given path should be an absolute path to a directory. If it's not
// a directory, the function will not proceed and return an empty list.
//
// This is a helper function for Provider implementations.
func execsInPath(path string) ([]string, error) {
	if !isDir(path) {
		return nil, nil
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var execs []string

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if !looksLikePython(entry.Name()) {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			return nil, err
		}
		if !isExecutable(info) {
			continue
		}
		execs = append(execs, filepath.Join(path, entry.Name()))
	}

	return execs, nil
}

// looksLikePython returns true if the given filename looks like a Python
// executable.
func looksLikePython(name string) bool {
	return pythonFileRegex.MatchString(name)
}

// isExecutable returns true if the given file info is executable.
// On Windows, it just checks if the file extension is ".exe" or not.
func isExecutable(info fs.FileInfo) bool {
	if runtime.GOOS == "windows" {
		return strings.ToLower(filepath.Ext(info.Name())) == ".exe"
	}
	return info.Mode().IsRegular() && info.Mode()&0o111 != 0
}
