package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pyvenv",
	Short: "Personal tool to manage Python virtual environments",
	Long: `A personal tool to manage Python virtual environments across
different OS using the builtin 'venv' module.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
