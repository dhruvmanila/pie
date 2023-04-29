package venv

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/dhruvmanila/pie/internal/xdg"
)

// Names returns the names of all the managed virtual environments in a
// sorted order. The order is determined by [os.ReadDir].
func Names() ([]string, error) {
	entries, err := os.ReadDir(xdg.DataDir)
	if err != nil {
		return nil, err
	}

	var venvs []string
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		venvs = append(venvs, entry.Name())
	}

	return venvs, nil
}

// ProjectPath returns the absolute path to the project this virtual
// environment belongs to. This information is extracted from the
// `.project` file present in the virtual environment directory.
func ProjectPath(venvDir string) (string, error) {
	content, err := os.ReadFile(filepath.Join(venvDir, ".project"))
	if err != nil {
		return "", err
	}
	content = bytes.TrimRight(content, "\r\n")
	return string(content), nil
}

// PythonVersion returns the Python version this environment was created from.
// This information is extracted from the config file present in the virtual
// environment directory.
func PythonVersion(venvDir string) (string, error) {
	file, err := os.Open(filepath.Join(venvDir, "pyvenv.cfg"))
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		key, value, found := strings.Cut(scanner.Text(), "=")
		if !found {
			continue
		}
		if strings.TrimSpace(key) == "version" {
			return strings.TrimSpace(value), nil
		}
	}

	return "", errors.New("venv config file does not contain 'version' key")
}
