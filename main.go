package main

import (
	"jn/prompts"
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

	// arguments := os.Args[1:]
}
