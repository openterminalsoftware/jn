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

	vaultPath := resolveVaultPath(config)
	oldState, err := utils.EnableRawMode(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println("Failed to enable raw mode:", err)
		return
	}
	defer utils.DisableRawMode(int(os.Stdin.Fd()), oldState)

	fmt.Print(colors.Bold + "Search: " + colors.Reset)
	var query string
	var results []SearchResult

	for {
		char, _, err := utils.ReadChar()
		if err != nil {
			fmt.Println("Error reading char:", err)
			break
		}
		switch char {
		case 13: // Enter
			return
		case 9: // TAB
			if len(results) > 0 {
				PreviewFile(results[0].FilePath)
			}
			return
		case 127: // Backspace
			if len(query) > 0 {
				query = query[:len(query)-1]
			}
		default:
			query += string(char)
		}

		fmt.Print("\033[H\033[2J")
		fmt.Printf(colors.Bold+"Search: "+colors.Reset+"%s\033[K\n\n", query)
		if len(query) > 0 {
			results = findMatchingFiles(vaultPath, query)
			fmt.Print(colors.Blue + "--- Results ---" + colors.Reset + "\033[K\n")
			if len(results) == 0 {
				fmt.Print(colors.Red + "No matches found." + colors.Reset + "\033[K\n")
			} else {
				for i, res := range results {
					fn := strings.TrimPrefix(res.FilePath, vaultPath+string(os.PathSeparator))
					fmt.Printf(colors.Bold+"%2d. "+colors.Reset+colors.Magenta+"%s"+colors.Reset+"\033[K\n", i+1, fn)
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
}

func resolveVaultPath(config map[string]interface{}) string {
	var vp string
	if config != nil {
		if raw, ok := config["vault"]; ok {
			vp, _ = raw.(string)
		}
	}
	if vp == "" {
		vp = os.ExpandEnv("$HOME/.jn/vault")
	} else if strings.HasPrefix(vp, "~/") {
		if hd, err := os.UserHomeDir(); err == nil {
			vp = filepath.Join(hd, vp[2:])
		}
	} else {
		vp = os.ExpandEnv(vp)
	}
	return vp
}

func findMatchingFiles(vaultPath, query string) []SearchResult {
	var matches []SearchResult
	lq := strings.ToLower(query)
	filepath.Walk(vaultPath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(info.Name(), ".md") {
			return err
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		s := string(content)
		if strings.Contains(strings.ToLower(info.Name()), lq) ||
			strings.Contains(strings.ToLower(s), lq) {
			snippet := ""
			for _, line := range strings.Split(s, "\n") {
				if strings.Contains(strings.ToLower(line), lq) {
					snippet = strings.TrimSpace(line)
					break
				}
			}
			matches = append(matches, SearchResult{FilePath: path, Snippet: snippet})
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
	offset := 0

	for {
		width, height, err := getTerminalSize(os.Stdout.Fd())
		if err != nil {
			width, height = 80, 24
		}

		fmt.Print("\033[H\033[2J")
		for i := 0; i < height-2 && offset < len(lines); i++ {
			line := lines[offset]
			offset++
			if utf8.RuneCountInString(line) > width {
				line = string([]rune(line)[:width-3]) + "..."
			}
			fmt.Println(line)
		}

		if offset >= len(lines) {
			break
		}

		fmt.Print(colors.DarkGray + "-- More (TAB to scroll, any other key to quit) --" + colors.Reset + "\n")
		char, _, err := utils.ReadChar()
		if err != nil || char != 9 {
			break
		}
	}
}
