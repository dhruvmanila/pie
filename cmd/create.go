package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	"github.com/dhruvmanila/pyvenv/internal/pathutil"
	"github.com/dhruvmanila/pyvenv/internal/project"
	"github.com/dhruvmanila/pyvenv/internal/pythonfinder"
	"github.com/spf13/cobra"
)

// pythonVersion is the Python version to use for creating the virtual environment.
var pythonVersion string

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a virtual environment",
	Long: `Create a virtual environment for the current directory.

The environment will be created using the builtin 'venv' module. If the
'--python' flag is not specified, the default Python version will be used.
`,
	Args: cobra.NoArgs,
	Run: func(_ *cobra.Command, _ []string) {
		bold.Println("==> Creating a virtualenv for this project...")

		p, err := project.NewFromWd()
		if err != nil {
			log.Fatal(err)
		}

		if pathutil.IsDir(p.VenvDir) {
			log.Fatal(red.Sprintf("✘ Virtualenv already exists for this project: %s", p.Name))
		} else if errors.Is(err, fs.ErrNotExist) {
			if err = createVenv(p); err != nil {
				if errors.Is(err, pythonfinder.ErrVersionNotFound) {
					if pythonVersion != "" {
						log.Fatal(red.Sprintf("✘ Python version %s does not exist!", pythonVersion))
					} else {
						log.Fatal(red.Sprintf("✘ No Python version found!"))
					}
				} else if errors.Is(err, pythonfinder.ErrInvalidVersion) {
					log.Fatal(red.Sprintf("✘ Invalid Python version: %s (expected format: <major>.<minor>.<patch>)", pythonVersion))
				}
				log.Fatal(err)
			} else {
				// Associate project directory with the environment.
				if err = p.WriteProjectFile(); err != nil {
					log.Fatal(err)
				}

				green.Println("✔ Successfully created virtual environment!")
				fmt.Printf("Virtualenv location: %s\n", green.Sprint(p.VenvDir))
			}
		} else if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().StringVar(
		&pythonVersion, "python", "", `specify which version of Python to use for
creating the virtualenv`,
	)
}

func createVenv(p *project.Project) error {
	v, err := pythonfinder.New().Find(pythonVersion)
	if err != nil {
		return err
	}

	fmt.Printf("Using %s %s to create virtualenv...\n",
		yellowBold.Sprint(v.Executable),
		green.Sprintf("(%s)", v.VersionInfo),
	)

	// Creating the virtual environment using the 'venv' module does not
	// produce any output, so the only output we get is the error message,
	// if the command fails, to stderr.
	var stderr bytes.Buffer

	cmd := exec.Command(v.Executable, "-m", "venv", p.VenvDir, "--prompt", p.Name)
	cmd.Stderr = &stderr

	if err = cmd.Start(); err != nil {
		return commandError(err, cmd, stderr)
	}

	// stop channel is used to signal that the command has finished and
	// to stop listening for signals.
	stop := make(chan struct{})

	// signalReceived channel will receive true if a signal was received
	// during the command execution, and false otherwise.
	signalReceived := make(chan bool, 1)
	defer close(signalReceived)

	// Start listening for signals.
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
		defer signal.Stop(sig)

		select {
		case <-sig:
			signalReceived <- true
		case <-stop:
			signalReceived <- false
		}
	}()

	err = cmd.Wait()

	// Stop listening for signals. Now, if a signal was sent, then the
	// goroutine has already stopped because the signalReceived channel
	// is buffered.
	close(stop)
	if <-signalReceived {
		// Ensure that the virtual environment is deleted if we received
		// a signal to cancel the command.
		os.RemoveAll(p.VenvDir)
		os.Exit(1)
	}

	// There was no signal received, so we can safely check the error.
	if err != nil {
		return commandError(err, cmd, stderr)
	}

	return nil
}

// commandError returns a formatted error message on command failure.
// The out parameter is the output of the command, which might contain
// useful information about the error.
func commandError(err error, cmd *exec.Cmd, out bytes.Buffer) error {
	var output string
	if out.Len() > 0 {
		output = fmt.Sprintf("Output:\n  %s",
			faint.Sprint(strings.ReplaceAll(out.String(), "\n", "\n  ")),
		)
	}

	return fmt.Errorf("Command failure: %s\n  %s\n%s",
		red.Sprint(err),
		bold.Sprint(strings.Join(cmd.Args, " ")),
		output,
	)
}
