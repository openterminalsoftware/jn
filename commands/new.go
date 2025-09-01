package new

import (
	"bufio"
	"fmt"
	"jn/colors"
	"jn/utils"
	"os"
	"path/filepath"
	"strings"
)

func New() {
	content := utils.TextEditor(colors.DarkGray + "Type your markdown \".exit\" to save & quit" + colors.Reset + "\r\n")
	if strings.TrimSpace(content) == "" {
		fmt.Println("\nAborted, empty content.")
		return
	}

	var config map[string]interface{}
	configPath := os.ExpandEnv("$HOME/.jn/config.json")
	if _, err := os.Stat(configPath); err == nil {
		config = utils.ParseConfig(configPath)
	}

	var vaultPath string
	if config != nil {
		vaultPathIface, ok := config["vault"]
		if ok {
			vaultPath, _ = vaultPathIface.(string)
		}
	}

	if vaultPath == "" {
		vaultPath = os.ExpandEnv("$HOME/.jn/vault") // Default value
	} else if strings.HasPrefix(vaultPath, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Println("Failed to get home directory:", err)
			return
		}
		vaultPath = filepath.Join(homeDir, vaultPath[2:])
	} else if strings.HasPrefix(vaultPath, "$HOME/") {
		vaultPath = os.ExpandEnv(vaultPath)
	}

	if err := os.MkdirAll(vaultPath, 0755); err != nil {
		fmt.Println("Failed to create vault directory:", err)
		return
	}

	fmt.Print("\nEnter a filename for your note: ")
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			fmt.Printf("\nError reading filename: %v. Aborting.\n", err)
		} else {
			fmt.Println("\nNo filename entered. Aborting.")
		}
		return
	}
	fileName := scanner.Text()

	if strings.TrimSpace(fileName) == "" {
		fmt.Println("\nFilename cannot be empty. Aborting.")
		return
	}

	if !strings.HasSuffix(strings.ToLower(fileName), ".md") {
		fileName = fileName + ".md"
	}

	fullPath := filepath.Join(vaultPath, fileName)
	if _, err := os.Stat(fullPath); err == nil {
		base := strings.TrimSuffix(fileName, ".md")
		for i := 1; ; i++ {
			newName := fmt.Sprintf("%s-%d.md", base, i)
			fullPath = filepath.Join(vaultPath, newName)
			if _, err := os.Stat(fullPath); os.IsNotExist(err) {
				break
			}
		}
	}

	err := os.WriteFile(fullPath, []byte(content), 0644)
	if err != nil {
		fmt.Println("Failed to save:", err)
		return
	}
	fmt.Println("\nEntry saved at:", fullPath)
}
