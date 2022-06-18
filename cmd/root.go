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
	noColor        bool
	outputVenvInfo bool
)

var rootCmd = &cobra.Command{
	Use:     "pyvenv",
	Short:   "Personal tool to manage Python virtual environments.",
	Version: Version,
	Run: func(_ *cobra.Command, _ []string) {
		if outputVenvInfo {
			p, err := project.NewFromWd()
			if err != nil {
				log.Fatal(err)
			}

			if stat, err := os.Stat(p.VenvDir); err == nil && stat.IsDir() {
				fmt.Println(p.VenvDir)
			}
		}
	},
}

func Execute() {
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
