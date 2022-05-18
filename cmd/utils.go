package cmd

import (
	"os"
	"path/filepath"
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
