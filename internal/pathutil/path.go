package pathutil

import "os"

// IsDir returns true if the given path exists and is a directory.
// It delegates to os.Stat and FileInfo.IsDir.
func IsDir(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}
