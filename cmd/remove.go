package cmd

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
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
	ValidArgsFunction: func(_ *cobra.Command, args []string, _ string) ([]string, cobra.ShellCompDirective) {
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		venvInfo, err := getVenvInfo()
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		return venvInfo.Names, cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(_ *cobra.Command, args []string) {
		dataDir, err := xdg.DataFile("pyvenv/")
		if err != nil {
			log.Fatal(err)
		}

		venvName, err := getVenvNameFromArgs(args)
		if err != nil {
			log.Fatal(err)
		}

		venvDir := filepath.Join(dataDir, venvName)
		fmt.Printf("Virtualenv location: %s\n", green.Sprint(venvDir))
		if stat, err := os.Stat(venvDir); err == nil && stat.IsDir() {
			if err = os.RemoveAll(venvDir); err != nil {
				log.Fatal(err)
			} else {
				green.Println("✔ Successfully removed virtual environment!")
			}
		} else if errors.Is(err, fs.ErrNotExist) {
			log.Fatal(red.Sprintf("✘ Virtualenv %q does not exist!", venvName))
		} else if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
