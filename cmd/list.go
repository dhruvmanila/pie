package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/dhruvmanila/pyvenv/internal/venv"
	"github.com/dhruvmanila/pyvenv/internal/xdg"
)

// verbose is a flag used to output additional environment information.
var verbose bool

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List out all the managed virtualenvs",
	Aliases: []string{"ls"},
	Args:    cobra.NoArgs,
	Run: func(_ *cobra.Command, _ []string) {
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
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "output additional venv information")
}
