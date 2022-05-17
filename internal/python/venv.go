package python

import (
	"fmt"
	"io"
	"os/exec"
	"path/filepath"

	"github.com/adrg/xdg"
)

// executeCommand will execute the given command returning any error faced
// during the execution.
func executeCommand(cmd *exec.Cmd) error {
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	if err = cmd.Start(); err != nil {
		return err
	}

	stderrOut, err := io.ReadAll(stderr)
	if err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("%s: %s", err, stderrOut)
	}

	return nil
}

// CreateVenv is used to create the virtual environment for the given
// Python version with the given name.
//
// This will store all the virtual environments in `XDG_DATA_HOME` directory.
// It will use the builtin `venv` module to create the virtual environment.
func CreateVenv(version string, name string) error {
	pythonExec, err := VersionLookup(version)
	if err != nil {
		return err
	}

	dir, err := xdg.DataFile("pyvenv/")
	if err != nil {
		return err
	}

	venvDir := filepath.Join(dir, name)
	cmd := exec.Command(pythonExec, "-m", "venv", venvDir)
	if err := executeCommand(cmd); err != nil {
		return err
	}

	return nil
}
