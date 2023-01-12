package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/dhruvmanila/pyvenv/internal/pythonfinder"
	"github.com/dhruvmanila/pyvenv/internal/venv"
	"github.com/dhruvmanila/pyvenv/internal/xdg"
)

var (
	// verbose is a flag used to output additional environment information.
	verbose bool

	// execs is a flag used to output all the available Python versions and
	// their executable paths.
	execs bool
)

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List out all the managed virtualenvs",
	Aliases: []string{"ls"},
	Args:    cobra.NoArgs,
	Run: func(_ *cobra.Command, _ []string) {
		if execs {
			printPythonVersions()
		} else {
			printVenvs()
		}
	},
}

func printPythonVersions() {
	versions, err := pythonfinder.New().FindAll()
	if err != nil {
		if errors.Is(err, pythonfinder.ErrVersionNotFound) {
			log.Fatal(red.Sprint("✘ No Python version found on the system"))
		}
		log.Fatal(err)
	}
	bold.Println("Found Python versions:")
	for _, v := range versions {
		fmt.Printf("  %s %s\n",
			yellowBold.Sprint(v.VersionInfo),
			faint.Sprintf("(%s)", v.Executable),
		)
	}
}

func printVenvs() {
	fmt.Printf("Root directory: %s\n", green.Sprint(xdg.DataDir))

	venvNames, err := venv.Names()
	if err != nil {
		log.Fatal(err)
	}

	_, currentVenvName := filepath.Split(os.Getenv("VIRTUAL_ENV"))
	for _, venvName := range venvNames {
		var line string
		if venvName == currentVenvName {
			line += bold.Sprint("* " + venvName)
		} else {
			line += bold.Sprint("  " + venvName)
		}
		if verbose {
			venvDir := filepath.Join(xdg.DataDir, venvName)
			projectPath, err := venv.ProjectPath(venvDir)
			if err != nil {
				log.Fatal(err)
			}
			pythonVersion, err := venv.PythonVersion(venvDir)
			if err != nil {
				log.Fatal(err)
			}
			line += yellowBold.Sprintf(" (%s)", pythonVersion) + faint.Sprintf(" (%s)", projectPath)
		}
		fmt.Println(line)
	}
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "output additional venv information")
	listCmd.Flags().BoolVar(&execs, "execs", false, "output available Python versions")
}
