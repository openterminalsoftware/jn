package utils

import (
	"fmt"
	"jn/colors"
	"os"
	"os/signal"
	"strings"
)

func TextEditor(prompt string) string {
	fd := int(os.Stdin.Fd())
	oldState, err := enableRawMode(fd)
	if err != nil {
		fmt.Println("Failed to enable raw mode:", err)
		os.Exit(1)
	}
	defer disableRawMode(fd, oldState)

	// Handle Ctrl+C clean exit
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		disableRawMode(fd, oldState)
		os.Exit(0)
	}()

	var content []string
	var line []rune

	buf := make([]byte, 1)
	fmt.Print(prompt)

	for {
		_, err := os.Stdin.Read(buf)
		if err != nil {
			break
		}
		ch := buf[0]

		if ch == 3 { // Ctrl+C
			break
		} else if ch == 13 { // Enter
			strLine := string(line)
			if strLine == ".exit" {
				break
			}
			content = append(content, strLine)
			fmt.Print("\r\n")
			line = []rune{}
		} else if ch == 127 { // Backspace
			if len(line) > 0 {
				line = line[:len(line)-1]
				fmt.Print("\b \b")
			}
		} else {
			line = append(line, rune(ch))
		}

		// redraw current line
		fmt.Print("\r\033[K" + colors.DarkGray + "| " + colors.Reset + HighlightMarkdown(string(line)))
	}

	return strings.Join(content, "\n")
}
