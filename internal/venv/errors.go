package venv

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// CommandError is an error that occurs when a command fails.
type CommandError struct {
	// Cmd is the command that failed.
	Cmd *exec.Cmd

	// Err is the error that occurred.
	Err error

	// Stderr is the standard error output of the command. It is only set if
	// the error is of type *exec.ExitError. Otherwise, it is nil.
	Stderr []byte
}

func (e *CommandError) Error() string {
	s := fmt.Sprintf("%q: %s", strings.Join(e.Cmd.Args, " "), e.Err)
	if e.Stderr != nil && len(e.Stderr) != 0 {
		s += fmt.Sprintf("\nOutput:\n  %s",
			bytes.ReplaceAll(e.Stderr, []byte("\n"), []byte("\n  ")),
		)
	}
	return s
}

func (e *CommandError) Unwrap() error {
	return e.Err
}
