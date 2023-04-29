package venv

import (
	"os/exec"

	"github.com/dhruvmanila/pie/internal/project"
	"github.com/dhruvmanila/pie/internal/pythonfinder"
)

// Create creates a new virtual environment for the given project using the
// given Python version.
//
// The virtual environment is created using the standard library's venv module.
// If the command execution fails, a *CommandError is returned. Other types
// of errors are returned as is.
//
// If the virtual environment gets created successfully, it will associate the
// project with the virtual environment. This is done by creating a '.project'
// file in the virtual environment directory which contains the absolute path
// to the project.
func Create(p *project.Project, v *pythonfinder.PythonVersion) error {
	cmd := exec.Command(v.Executable, "-m", "venv", p.VenvDir, "--prompt", p.Name)
	_, err := cmd.Output()
	if err != nil {
		cmdErr := &CommandError{Cmd: cmd, Err: err}
		if exitErr, ok := err.(*exec.ExitError); ok {
			cmdErr.Stderr = exitErr.Stderr
		}
		return cmdErr
	}
	// Associate project directory with the environment.
	return p.WriteProjectFile()
}
