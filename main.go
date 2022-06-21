package main

import (
	"log"

	"github.com/dhruvmanila/pyvenv/cmd"
)

// version is the current tool version. This is populated by goreleaser while
// building the binaries.
var version = "dev"

func main() {
	log.SetFlags(0)

	cmd.Execute(version)
}
