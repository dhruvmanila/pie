package cmd

import (
	"bytes"
	"os"
	"path/filepath"
)

// readProjectFile reads and returns the content of the project file
// present in the given virtual environment directory.
func readProjectFile(venvDir string) (string, error) {
	projectFilePath := filepath.Join(venvDir, ".project")
	content, err := os.ReadFile(projectFilePath)
	if err != nil {
		return "", err
	}
	content = bytes.TrimRight(content, "\n")
	return string(content), nil
}
