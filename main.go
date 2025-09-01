package main

import (
	"fmt"
	"os"
)

func help() {
	fmt.Println("jn. ")
}

func main() {
	if len(os.Args) < 2 {
		help()
		os.Exit(1)
	}

	// arguments := os.Args[1:]
}
