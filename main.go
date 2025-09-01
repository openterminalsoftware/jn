package main

import (
	"fmt"
	"os"
)

const (
	Red      = "\033[31m"
	Green    = "\033[32m"
	Yellow   = "\033[33m"
	Blue     = "\033[34m"
	Magenta  = "\033[35m"
	Cyan     = "\033[36m"
	Reset    = "\033[0m"
	Bold     = "\033[1m"
	Italic   = "\033[3m"
	DarkGray = "\033[90m"
)

func help() {
	fmt.Println(Bold + Magenta + "jn" + Reset)
	fmt.Println(Bold + Blue + "— Markdown previewer and journal/note-taking command line tool written in Go. " + Reset)
	fmt.Println("\nUsage:")
	fmt.Println(Bold + "  jn " + Reset + Italic + "<command>" + " [arguments]" + Reset)
	fmt.Println(Bold + "  jn " + Reset + Italic + "<file>" + Reset)
	fmt.Println("\nCommands:")
	fmt.Println(Bold + Cyan + "  help" + Reset + Green + "            Alias: h" + Reset)
	fmt.Println(Bold + Cyan + "  new" + Reset + Green + "             Alias: n" + Reset)
	fmt.Println(Bold + Cyan + "  preview" + Reset + Green + "         Alias: p" + Reset)
	fmt.Println(Bold + Cyan + "  publish" + Reset + Green + "         Alias: pub" + Reset)
	fmt.Println(Bold + Cyan + "  delete" + Reset + Green + "          Alias: d" + Reset)
	fmt.Println(Bold + Cyan + "  list" + Reset + Green + "            Alias: l" + Reset)
	fmt.Println(Bold + Cyan + "  search" + Reset + Green + "          Alias: s" + Reset)
	fmt.Println(Bold + Cyan + "  version" + Reset + Green + "         Alias: v" + Reset)
	fmt.Println(Bold + Cyan + "  config" + Reset + Green + "          Alias: conf" + Reset)

	fmt.Println("\n\n" + DarkGray + "jn help <command> for more information on a command." + Reset)

}

func main() {
	if len(os.Args) < 2 {
		help()
		os.Exit(1)
	} else if os.Args[1] == "help" || os.Args[1] == "h" {
		help()
		os.Exit(0)
	}

	// arguments := os.Args[1:]
}
