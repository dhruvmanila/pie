package cmd

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"strings"

	"github.com/dhruvmanila/pyvenv/internal/pathutil"
	"github.com/dhruvmanila/pyvenv/internal/project"
	"github.com/spf13/cobra"
)

// noConfirm is a flag to skip the confirmation prompt for removing a
// virtual environment.
var noConfirm bool

var removeCmd = &cobra.Command{
	Use:     "remove",
	Short:   "Remove the virtual environment",
	Long:    `Remove the virtual environment associated with the current project.`,
	Aliases: []string{"rm"},
	Args:    cobra.NoArgs,
	Run: func(_ *cobra.Command, _ []string) {
		p, err := project.NewFromWd()
		if err != nil {
			log.Fatal(err)
		}

		if pathutil.IsDir(p.VenvDir) {
			fmt.Printf("Removing virtualenv (%s)...\n", green.Sprint(p.VenvDir))

			var response string
			if !noConfirm {
				fmt.Print("Proceed? [y/n]: ")
				if _, err = fmt.Scan(&response); err != nil {
					log.Fatal(err)
				}
			}

			if noConfirm || strings.ToLower(strings.TrimSpace(response)) == "y" {
				if err = os.RemoveAll(p.VenvDir); err != nil {
					log.Fatal(err)
				} else {
					green.Println("✔ Successfully removed virtual environment!")
				}
			}
		} else if errors.Is(err, fs.ErrNotExist) {
			log.Fatal(red.Sprint("✘ No virtualenv has been created for this project yet!"))
		} else if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	removeCmd.Flags().BoolVarP(&noConfirm, "yes", "y", false, "skip the confirmation prompt")
	rootCmd.AddCommand(removeCmd)
}
