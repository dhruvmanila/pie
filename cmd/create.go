package cmd

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"

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

		if stat, err := os.Stat(p.VenvDir); err == nil && stat.IsDir() {
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

	cmd := exec.Command(v.Executable, "-m", "venv", p.VenvDir, "--prompt", p.Name)

	// Creating the virtual environment using the 'venv' module does not
	// produce any output, so the only output we get is the error message
	// if the command fails.
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %s", err, output)
	}

	return nil
}
