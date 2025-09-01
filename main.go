package main

import (
	"fmt"
	"os"
)

const (
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	Reset   = "\033[0m"
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
