package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

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
		fmt.Printf("Virtualenvs location: %s\n", green.Sprint(dataDir))

		entries, err := os.ReadDir(dataDir)
		if err != nil {
			log.Fatal(err)
		}

		_, currentVenv := filepath.Split(os.Getenv("VIRTUAL_ENV"))
		for _, entry := range entries {
			if entry.IsDir() {
				venvName := entry.Name()
				if currentVenv == venvName {
					bold.Println("* " + venvName)
				} else {
					bold.Println("  " + venvName)
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
