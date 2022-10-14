package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/dhruvmanila/pyvenv/internal/project"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	noColor        bool
	outputVenvInfo bool
)

var rootCmd = &cobra.Command{
	Use:   "pyvenv",
	Short: "Personal tool to manage Python virtual environments.",
	Run: func(_ *cobra.Command, _ []string) {
		if outputVenvInfo {
			path, err := os.Getwd()
			if err != nil {
				log.Fatal(err)
			}

			path, err = filepath.EvalSymlinks(path)
			if err != nil {
				log.Fatal(err)
			}

			var root string
			if runtime.GOOS == "windows" {
				// Windows root path - "C:" + "\"
				root = filepath.VolumeName(path) + string(os.PathSeparator)
			} else {
				root = "/"
			}

			for path != root {
				p, err := project.New(path)
				if err != nil {
					log.Fatal(err)
				}

				if stat, err := os.Stat(p.VenvDir); err == nil && stat.IsDir() {
					fmt.Println(p.VenvDir)
					break
				}

				path = filepath.Dir(path)
			}
		}
	},
}

func Execute(version string) {
	rootCmd.Version = version
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(setColorOutput)
	rootCmd.Flags().BoolVar(&outputVenvInfo, "venv", false, "output virtualenv information")
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "disable color output")
}

func setColorOutput() {
	if noColor {
		color.NoColor = true
	}
}
