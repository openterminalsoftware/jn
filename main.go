package main

import (
	"fmt"
	"jn/commands"
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
		commands.New()
	} else if utils.Contains([]string{"preview", "p"}, command) {
		commands.Preview(os.Args[2])
	} else if utils.Contains([]string{"delete", "d"}, command) {
		if len(os.Args) < 3 {
			prompts.Help()
			os.Exit(1)
		}
		commands.Delete(os.Args[2])
	} else if utils.Contains([]string{"list", "l"}, command) {
		fmt.Println(commands.List())
	} else {
		prompts.Help()
		os.Exit(1)
	}
}
