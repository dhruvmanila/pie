package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/dhruvmanila/pie/internal/pathutil"
	"github.com/dhruvmanila/pie/internal/venv"
	"github.com/dhruvmanila/pie/internal/xdg"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Remove any dangling virtual environments",
	Long: `Remove any dangling virtual environments.

A dangling virtual environment is one that is not associated with a project.
`,
	Args: cobra.NoArgs,
	Run: func(_ *cobra.Command, _ []string) {
		venvNames, err := venv.Names()
		if err != nil {
			log.Fatal(err)
		}

		count := 0
		for _, venvName := range venvNames {
			projectPath, err := venv.ProjectPath(venvName)
			if err != nil {
				log.Fatal(err)
			}

			if !pathutil.IsDir(projectPath) {
				venvDir := filepath.Join(xdg.DataDir, venvName)
				fmt.Printf("Removing virtualenv (%s)...\n", green.Sprint(venvDir))
				if err = os.RemoveAll(venvDir); err != nil {
					log.Fatal(err)
				}
				count++
			}
		}

		if count == 0 {
			green.Println("✔ No dangling virtual environments found")
		} else {
			green.Printf("✔ Removed %d dangling virtual environments\n", count)
		}
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)
}
