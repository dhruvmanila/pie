package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/adrg/xdg"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List out all the managed venvs",
	Aliases: []string{"ls"},
	Args:    cobra.NoArgs,
	Run: func(_ *cobra.Command, _ []string) {
		dataDir, err := xdg.DataFile("pyvenv/")
		if err != nil {
			log.Fatal(err)
		}

		entries, err := os.ReadDir(dataDir)
		if err != nil {
			log.Fatal(err)
		}

		for _, entry := range entries {
			if entry.IsDir() {
				fmt.Println(entry.Name())
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
