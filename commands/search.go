package commands

import (
	"fmt"
	"jn/utils"
	"os"
	"path/filepath"
	"strings"
)

type SearchResult struct {
	FilePath string
	Snippet  string
}

func Search() {
	var config map[string]interface{}
	configPath := os.ExpandEnv("$HOME/.jn/config.json")
	if _, err := os.Stat(configPath); err == nil {
		config = utils.ParseConfig(configPath)
	}

	var vaultPath string
	if config != nil {
		if vaultPathIface, ok := config["vault"]; ok {
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

	// Ensure vaultPath exists and is expanded
	vaultPath = os.ExpandEnv(vaultPath)
	if strings.HasPrefix(vaultPath, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Println("Failed to get home directory:", err)
			return
		}
		vaultPath = filepath.Join(homeDir, vaultPath[2:])
	}

	// 2. Terminal Raw Mode
	oldState, err := utils.EnableRawMode(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println("Failed to enable raw mode:", err)
		return
	}
	defer utils.DisableRawMode(int(os.Stdin.Fd()), oldState)

	fmt.Print("Search: ") // Initial prompt

	var query string
	var results []SearchResult

	for {
		char, _, err := utils.ReadChar()
		if err != nil {
			fmt.Println("Error reading char:", err)
			break
		}

		if char == 13 { // Enter
			break
		} else if char == 9 { // TAB
			if len(results) > 0 {
				query = ""
				Preview(results[0].FilePath)
				break
			}
		} else if char == 127 { // Backspace
			if len(query) > 0 {
				query = query[:len(query)-1]
			}
		} else {
			query += string(char)
		}

		fmt.Print("\033[H\033[2J")

		fmt.Printf("Search: %s\033[K\n", query)

		fmt.Print("\033[3;1H")

		if len(query) > 0 {
			results = findMatchingFiles(vaultPath, query)
			fmt.Println("--- Results ---\033[K")
			if len(results) == 0 {
				fmt.Println("No matches found.\033[K")
			} else {
				for i, res := range results {
					fmt.Printf("%2d. %s\033[K\n",
						i+1,
						strings.TrimPrefix(res.FilePath, vaultPath+string(os.PathSeparator)))
					if res.Snippet != "" {
						fmt.Printf("    %s\033[K\n", res.Snippet)
					}
				}
			}
		} else {
			fmt.Println("--- Results ---\033[K")
			fmt.Println("Type to search...\033[K")
		}
	}

	fmt.Printf("\nFinal search query: %s\n", query)
}

func findMatchingFiles(vaultPath, query string) []SearchResult {
	var matches []SearchResult
	lowerQuery := strings.ToLower(query)

	filepath.Walk(vaultPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(info.Name(), ".md") {
			return nil
		}

		foundInFile := false

		if strings.Contains(strings.ToLower(info.Name()), lowerQuery) {
			matches = append(matches, SearchResult{
				FilePath: path,
				Snippet:  "",
			})
			foundInFile = true
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		lines := strings.Split(string(content), "\n")
		for _, line := range lines {
			if strings.Contains(strings.ToLower(line), lowerQuery) {
				if !foundInFile {
					matches = append(matches, SearchResult{
						FilePath: path,
						Snippet:  strings.TrimLeft(line, " \t"),
					})
					foundInFile = true
				}
				break
			}
		}
		return nil
	})
	return matches
}
