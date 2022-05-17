package cmd

import (
	"errors"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove [flags] [name]",
	Short: "Remove the venv for the current directory",
	Long: `Remove the virtual environment for the current directory.

If the name of the environment is given, then the virtual environment
associated with that name is removed instead.
`,
	Aliases: []string{"rm"},
	Args:    cobra.MaximumNArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		dataDir := dataDir()

		venvName, err := getVenvNameFromArgs(args)
		if err != nil {
			log.Fatal(err)
		}

		venvDir := filepath.Join(dataDir, venvName)
		if stat, err := os.Stat(venvDir); err == nil && stat.IsDir() {
			if err = os.RemoveAll(venvDir); err != nil {
				log.Fatal(err)
			}
		} else if errors.Is(err, fs.ErrNotExist) {
			log.Fatalf("%s: venv does not exist", venvName)
		} else if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
