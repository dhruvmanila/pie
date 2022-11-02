package pythonfinder

import (
	"os"
)

// isDir returns true if the given path exists and is a directory.
// It delegates to os.Stat and FileInfo.IsDir.
func isDir(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}
