package commands

import (
	"fmt"
	"jn/colors"
	"jn/utils"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"unicode/utf8"
	"unsafe"
)

type winsize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

func getTerminalSize(fd uintptr) (width, height int, err error) {
	ws := winsize{}
	_, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		fd,
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(&ws)),
	)
	if errno != 0 {
		err = errno
		return
	}
	width = int(ws.Col)
	height = int(ws.Row)
	return
}

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

	// Expand any env in vaultPath
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

	fmt.Print(colors.Bold + "Search: " + colors.Reset) // Initial prompt

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
				PreviewFile(results[0].FilePath) // show full content
				break
			}
		} else if char == 127 { // Backspace
			if len(query) > 0 {
				query = query[:len(query)-1]
			}
		} else {
			query += string(char)
		}

		// Clear screen
		fmt.Print("\033[H\033[2J")

		// Print search prompt
		fmt.Printf(colors.Bold+"Search: "+colors.Reset+"%s\033[K\n\n", query)

		if len(query) > 0 {
			results = findMatchingFiles(vaultPath, query)
			fmt.Print(colors.Blue + "--- Results ---" + colors.Reset + "\033[K\n")
			if len(results) == 0 {
				fmt.Print(colors.Red + "No matches found." + colors.Reset + "\033[K\n")
			} else {
				for i, res := range results {
					fileName := strings.TrimPrefix(res.FilePath, vaultPath+string(os.PathSeparator))
					fmt.Printf(colors.Bold+"%2d. "+colors.Reset+colors.Magenta+"%s"+colors.Reset+"\033[K\n",
						i+1,
						fileName)
					if res.Snippet != "" {
						fmt.Printf("    %s\033[K\n", utils.HighlightMarkdown(res.Snippet))
					}
				}
			}
		} else {
			fmt.Print(colors.Blue + "--- Results ---" + colors.Reset + "\033[K\n")
			fmt.Print(colors.DarkGray + "Type to search..." + colors.Reset + "\033[K\n")
		}
	}

	fmt.Printf("\n"+colors.DarkGray+"Final search query: %s"+colors.Reset+"\n", query)
}

func findMatchingFiles(vaultPath, query string) []SearchResult {
	var matches []SearchResult
	lowerQuery := strings.ToLower(query)

	filepath.Walk(vaultPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(info.Name(), ".md") {
			return nil
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return nil // Continue walking
		}

		contentStr := string(content)
		lowerContent := strings.ToLower(contentStr)

		fileNameMatch := strings.Contains(strings.ToLower(info.Name()), lowerQuery)
		contentMatch := strings.Contains(lowerContent, lowerQuery)

		if fileNameMatch || contentMatch {
			snippet := ""

			// If content matched, grab first line containing query
			if contentMatch {
				lines := strings.Split(contentStr, "\n")
				for _, line := range lines {
					if strings.Contains(strings.ToLower(line), lowerQuery) {
						snippet = strings.TrimSpace(line)
						break
					}
				}
			}

			matches = append(matches, SearchResult{
				FilePath: path,
				Snippet:  snippet,
			})
		}
		return nil
	})
	return matches
}
func PreviewFile(path string) {
	content, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	lines := strings.Split(string(content), "\n")

	// Get terminal size
	width, height, err := getTerminalSize(os.Stdout.Fd())
	if err != nil {
		width, height = 80, 24 // fallback
	}

	for i, line := range lines {
		// Truncate overly long lines
		if utf8.RuneCountInString(line) > width {
			line = string([]rune(line)[:width-3]) + "..."
		}
		fmt.Println(line)

		// Paginate based on screen height
		if i+1 >= height-2 {
			fmt.Print(colors.DarkGray + "-- More (TAB to scroll, q to quit) --" + colors.Reset + "\n")
			break
		}
	}
}
