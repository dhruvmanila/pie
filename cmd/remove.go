package cmd

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"

	"github.com/dhruvmanila/pyvenv/internal/project"
	"github.com/spf13/cobra"
)

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

		if stat, err := os.Stat(p.VenvDir); err == nil && stat.IsDir() {
			fmt.Printf("Removing virtualenv (%s)...\n", green.Sprint(p.VenvDir))
			if err = os.RemoveAll(p.VenvDir); err != nil {
				log.Fatal(err)
			} else {
				green.Println("✔ Successfully removed virtual environment!")
			}
		} else if errors.Is(err, fs.ErrNotExist) {
			log.Fatal(red.Sprint("✘ No virtualenv has been created for this project yet!"))
		} else if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
