package main

import (
	"jn/commands/new"
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

	arguments := os.Args[1:]

	if utils.Contains([]string{"new", "n"}, arguments[0]) {
		new.New()
	} else {
		prompts.Help()
		os.Exit(1)
	}
}
