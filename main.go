package main

import (
	"fmt"
	"jn/colors"
	"os"
)

func help() {
	fmt.Println(colors.Bold + colors.Magenta + "jn" + colors.Reset)
	fmt.Println(colors.Bold + colors.Blue + "— Markdown previewer and journal/note-taking command line tool written in Go. " + colors.Reset)
	fmt.Println("\nUsage:")
	fmt.Println(colors.Bold + "  jn " + colors.Reset + colors.Italic + "<command>" + " [arguments]" + colors.Reset)
	fmt.Println(colors.Bold + "  jn " + colors.Reset + colors.Italic + "<file>" + colors.Reset)
	fmt.Println("\nCommands:")
	fmt.Println(colors.Bold + colors.Cyan + "  help" + colors.Reset + colors.Green + "            Alias: h" + colors.Reset)
	fmt.Println(colors.Bold + colors.Cyan + "  new" + colors.Reset + colors.Green + "             Alias: n" + colors.Reset)
	fmt.Println(colors.Bold + colors.Cyan + "  preview" + colors.Reset + colors.Green + "         Alias: p" + colors.Reset)
	fmt.Println(colors.Bold + colors.Cyan + "  publish" + colors.Reset + colors.Green + "         Alias: pub" + colors.Reset)
	fmt.Println(colors.Bold + colors.Cyan + "  delete" + colors.Reset + colors.Green + "          Alias: d" + colors.Reset)
	fmt.Println(colors.Bold + colors.Cyan + "  list" + colors.Reset + colors.Green + "            Alias: l" + colors.Reset)
	fmt.Println(colors.Bold + colors.Cyan + "  search" + colors.Reset + colors.Green + "          Alias: s" + colors.Reset)
	fmt.Println(colors.Bold + colors.Cyan + "  version" + colors.Reset + colors.Green + "         Alias: v" + colors.Reset)
	fmt.Println(colors.Bold + colors.Cyan + "  config" + colors.Reset + colors.Green + "          Alias: conf" + colors.Reset)

	fmt.Println("\n\n" + colors.DarkGray + "jn help <command> for more information on a command." + colors.Reset)

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
