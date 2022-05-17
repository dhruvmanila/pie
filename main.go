package main

import (
	"log"

	"github.com/dhruvmanila/pyvenv/cmd"
)

func main() {
	log.SetPrefix("pyvenv: ")
	log.SetFlags(0)

	cmd.Execute()
}
