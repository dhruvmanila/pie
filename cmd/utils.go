package cmd

import (
	"bytes"
	"os"
	"path/filepath"
)

func readProjectFile(venvDir string) (string, error) {
	projectFilePath := filepath.Join(venvDir, ".project")
	content, err := os.ReadFile(projectFilePath)
	if err != nil {
		return "", err
	}
	content = bytes.TrimRight(content, "\n")
	return string(content), nil
}
