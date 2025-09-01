package commands

import (
	"fmt"
	"io/ioutil"
	"jn/utils"
	"os"
	"strings"
)

func Preview(filePath string) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		fmt.Println(utils.HighlightMarkdown(line))
	}
}
