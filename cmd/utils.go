package cmd

import (
	"os"
	"path/filepath"
	"runtime"
)

// getVenvNameFromArgs returns the virtual environment name from the given
// list of arguments. If name is not provided in the arguments, then the
// last part of the current working directory is returned.
func getVenvNameFromArgs(args []string) (string, error) {
	if len(args) == 1 {
		return args[0], nil
	}

	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	_, venvName := filepath.Split(cwd)
	return venvName, nil
}

// dataDir returns the data directory for the tool.
//
// The precedences is as follows:
//   1. XDG_DATA_HOME
//   2. LocalAppData (windows only)
//   3. HOME
func dataDir() string {
	var path string
	if xdgDataHome := os.Getenv("XDG_DATA_HOME"); xdgDataHome != "" {
		path = filepath.Join(xdgDataHome, "pyvenv")
	} else if localAppData := os.Getenv("LocalAppData"); runtime.GOOS == "windows" && localAppData != "" {
		path = filepath.Join(localAppData, "pyvenv")
	} else {
		homeDir, _ := os.UserHomeDir()
		path = filepath.Join(homeDir, ".local", "share", "pyvenv")
	}
	return path
}
