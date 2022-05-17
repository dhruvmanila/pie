package cmd

import (
	"errors"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/dhruvmanila/pyvenv/internal/python"
	"github.com/spf13/cobra"
)

// pythonVersion is the Python version to use to create the virtual environment.
var pythonVersion string

var createCmd = &cobra.Command{
	Use:   "create [flags] [name]",
	Short: "Create a venv for the current directory",
	Long: `Create a virtual environment for the current directory.

The environment will be created using the builtin 'venv' module. If the
'--version' flag is not specified, the default Python version will be used.

The name argument is used to name the virtual environment. If it is not
provided, then the project name will be used.
`,
	Args: cobra.MaximumNArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		dataDir := dataDir()

		venvName, err := getVenvNameFromArgs(args)
		if err != nil {
			log.Fatal(err)
		}

		venvDir := filepath.Join(dataDir, venvName)
		if stat, err := os.Stat(venvDir); err == nil && stat.IsDir() {
			log.Fatalf("%s: venv exists for project", venvName)
		} else if errors.Is(err, fs.ErrNotExist) {
			if err = python.CreateVenv(pythonVersion, venvName); err != nil {
				log.Fatal(err)
			}
		} else if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().StringVarP(&pythonVersion, "version", "v", "", "Use this Python version instead")
}