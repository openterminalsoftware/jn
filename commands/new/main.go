package new

import (
	"bufio"
	"fmt"
	"jn/colors"
	"os"
	"path/filepath"
	"strings"
)

func New() {
	configPath := os.ExpandEnv("$HOME/.jn/config.json")
	configDir := filepath.Dir(configPath)
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		fmt.Println("Creating ~/.jn/config.json")
		if err := os.MkdirAll(configDir, 0755); err != nil {
			fmt.Printf("Failed to create directory: %v\n", err)
			return
		}
		f, err := os.Create(configPath)
		if err != nil {
			fmt.Printf("Failed to create config file: %v\n", err)
			return
		}
		f.Close()
	}

	fmt.Println(colors.DarkGray + "Using ~/.jn/config.json" + colors.Reset)
	makeMarkdownEntry()
}

func TextEditor(initialContent string) (string, error) {
	fmt.Println("Enter your markdown content below. Type '.exit' on a new line to save and exit.")
	fmt.Println(colors.DarkGray + "------------------------------------------------------------" + colors.Reset)
	if initialContent != "" {
		fmt.Println(initialContent)
	}
	var lines []string
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(colors.DarkGray + "| " + colors.Reset)
		line, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		line = strings.TrimRight(line, "\r\n")
		if line == ".exit" || line == ":wq" {
			break
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n"), nil
}

func makeMarkdownEntry() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Entry name: ")
	fileNameInput, _ := reader.ReadString('\n')
	fileName := strings.TrimSpace(fileNameInput)
	if fileName == "" {
		fmt.Println("File name cannot be empty")
		return
	}
	if !strings.HasSuffix(fileName, ".md") {
		fileName += ".md"
	}

	vaultDir := os.ExpandEnv("$HOME/.jn/vault")
	if err := os.MkdirAll(vaultDir, 0755); err != nil {
		fmt.Printf("Failed to create vault directory: %v\n", err)
		return
	}
	fullPath := filepath.Join(vaultDir, fileName)

	var initialContent string
	if _, err := os.Stat(fullPath); err == nil {
		data, err := os.ReadFile(fullPath)
		if err == nil {
			initialContent = string(data)
		}
	}

	content, err := TextEditor(initialContent)
	if err != nil {
		fmt.Printf("Error in editor: %v\n", err)
		return
	}

	// Write content to file
	err = os.WriteFile(fullPath, []byte(content), 0644)
	if err != nil {
		fmt.Printf("Failed to write markdown file: %v\n", err)
		return
	}

	fmt.Println("Entry created successfully:", fullPath)
}

func main() {
	New()
}
