package cmd

import "github.com/fatih/color"

var (
	bold       = color.New(color.Bold)
	yellowBold = bold.Add(color.FgYellow)
	green      = color.New(color.FgGreen)
	red        = color.New(color.FgRed)
	faint      = color.New(color.Faint)
)
