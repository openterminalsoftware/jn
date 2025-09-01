package main

import (
	"jn/commands/new"
	"jn/commands/preview"
	"jn/prompts"
	"jn/utils"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		prompts.Help()
		os.Exit(1)
	} else if os.Args[1] == "help" || os.Args[1] == "h" {
		prompts.Help()
		os.Exit(0)
	}

	// arguments := os.Args[2:]
	command := os.Args[1]

	if utils.Contains([]string{"new", "n"}, command) {
		new.New()
	} else if utils.Contains([]string{"preview", "p"}, command) {
		preview.Preview() // Initialize preview command
	} else {
		prompts.Help()
		os.Exit(1)
	}
}
