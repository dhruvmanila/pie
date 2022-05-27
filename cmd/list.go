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
		venvInfo, err := getVenvInfo()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Virtualenvs location: %s\n", green.Sprint(venvInfo.RootDir))

		_, currentVenv := filepath.Split(os.Getenv("VIRTUAL_ENV"))
		for _, venvName := range venvInfo.Names {
			if currentVenv == venvName {
				bold.Println("* " + venvName)
			} else {
				bold.Println("  " + venvName)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

// VenvInfo contains information regarding all the virtual environments
// managed by the tool.
type VenvInfo struct {
	// RootDir is the directory where all the virtual environments are stored.
	RootDir string

	// Names is a list of all the available virtual environments.
	Names []string
}

// getVenvInfo returns information regarding all the managed virtual environments.
func getVenvInfo() (*VenvInfo, error) {
	dataDir, err := xdg.DataFile("pyvenv/")
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(dataDir)
	if err != nil {
		return nil, err
	}

	venvNames := make([]string, 0)
	for _, entry := range entries {
		if entry.IsDir() {
			venvNames = append(venvNames, entry.Name())
		}
	}
	return &VenvInfo{
		RootDir: dataDir,
		Names:   venvNames,
	}, nil
}
