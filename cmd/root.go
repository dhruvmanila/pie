package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/dhruvmanila/pyvenv/internal/project"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	// noColor is a flag to disable color output.
	noColor bool

	// outputVenvInfo is a flag to output the absolute path to the
	// virtual environment for the current project if there is any.
	outputVenvInfo bool
)

var rootCmd = &cobra.Command{
	Use:   "pyvenv",
	Short: "A tool to manage Python virtual environments.",
	Run: func(_ *cobra.Command, _ []string) {
		if outputVenvInfo {
			p, err := project.Current()
			if err != nil {
				log.Fatal(err)
			}

			if p != nil {
				fmt.Println(p.VenvDir)
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
