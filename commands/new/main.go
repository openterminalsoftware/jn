package new

import (
	"bufio"
	"fmt"
	"jn/colors"
	"jn/utils"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"strings"
)

func highlightMarkdown(line string) string {
	if strings.HasPrefix(line, "#") {
		return colors.Blue + line + colors.Reset
	}

	boldAsteriskRe := regexp.MustCompile(`\*\*([^\*]+)\*\*`)
	line = boldAsteriskRe.ReplaceAllString(line, colors.Bold+"$1"+colors.Reset)
	boldUnderscoreRe := regexp.MustCompile(`__([^_]+)__`)
	line = boldUnderscoreRe.ReplaceAllString(line, colors.Bold+"$1"+colors.Reset)

	italicAsteriskRe := regexp.MustCompile(`\*([^\*]+)\*`)
	line = italicAsteriskRe.ReplaceAllString(line, colors.Italic+"$1"+colors.Reset)
	italicUnderscoreRe := regexp.MustCompile(`_([^_]+)_`)
	line = italicUnderscoreRe.ReplaceAllString(line, colors.Italic+"$1"+colors.Reset)

	codeRe := regexp.MustCompile("`([^`]+)`")
	line = codeRe.ReplaceAllString(line, colors.Cyan+"$1"+colors.Reset)

	linkRe := regexp.MustCompile(`\[(.*?)](.*?)`)
	line = linkRe.ReplaceAllString(line,
		colors.Magenta+"[$1]"+colors.Reset+
			"("+colors.Blue+"$2"+colors.Reset+")")

	if strings.HasPrefix(line, "- ") || strings.HasPrefix(line, "* ") || strings.HasPrefix(line, "+ ") {
		return colors.Green + line + colors.Reset
	}
	if strings.HasPrefix(line, ">") {
		return colors.Yellow + line + colors.Reset
	}
	if strings.HasPrefix(line, "```") {
		return colors.Cyan + line + colors.Reset
	}
	return line
}

func TextEditor() string {
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
	fmt.Print(colors.DarkGray + "Type your markdown \".exit\" to save & quit" + colors.Reset + "\r\n")

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
		fmt.Print("\r\033[K" + colors.DarkGray + "| " + colors.Reset + highlightMarkdown(string(line)))
	}

	return strings.Join(content, "\n")
}

// New is an importable function to create a new entry interactively.
func New() {
	content := TextEditor()
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
