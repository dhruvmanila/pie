package cmd

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"os/exec"

	"github.com/dhruvmanila/pyvenv/internal/project"
	"github.com/dhruvmanila/pyvenv/internal/python"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// pythonVersion is the Python version to use to create the virtual environment.
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

		p, err := project.New()
		if err != nil {
			log.Fatal(err)
		}

		if stat, err := os.Stat(p.VenvDir); err == nil && stat.IsDir() {
			log.Fatal(red.Sprintf("✘ Virtualenv already exists for this project: %s", p.Name))
		} else if errors.Is(err, fs.ErrNotExist) {
			if err = createVenv(p); err != nil {
				if _, ok := err.(*python.VersionNotFoundError); ok {
					log.Fatal(red.Sprintf("✘ Python version %s does not exist!", pythonVersion))
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
	pythonExecInfo, err := python.VersionLookup(pythonVersion)
	if err != nil {
		return err
	}

	fmt.Printf("Using %s %s to create virtualenv...\n",
		bold.Add(color.FgYellow).Sprint(pythonExecInfo.Path),
		green.Sprintf("(%s)", pythonExecInfo.Version),
	)

	cmd := exec.Command(pythonExecInfo.Path, "-m", "venv", p.VenvDir, "--prompt", p.Name)
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
