package commands

import (
	"fmt"
	"jn/colors"
	"os"
	"strings"
)

func List() string {
	vaultPath := "$HOME/.jn/vault"
	homeDir, err := os.UserHomeDir()
	if err == nil {
		vaultPath = strings.Replace(vaultPath, "$HOME", homeDir, 1)
	}

	var files []string
	var walk func(string)
	walk = func(dir string) {
		entries, err := os.ReadDir(dir)
		if err != nil {
			return
		}
		for _, entry := range entries {
			if entry.IsDir() {
				walk(dir + string(os.PathSeparator) + entry.Name())
				continue
			}
			if strings.HasSuffix(entry.Name(), ".md") {
				fullPath := dir + string(os.PathSeparator) + entry.Name()
				relPath := strings.TrimPrefix(fullPath, vaultPath)
				if strings.HasPrefix(relPath, string(os.PathSeparator)) {
					relPath = relPath[1:]
				}
				files = append(files, relPath)
			}
		}
	}

	walk(vaultPath)

	var sb strings.Builder
	sb.WriteString(colors.Blue + "Notes in vault:\n" + colors.Reset)
	for i, f := range files {
		sb.WriteString(fmt.Sprintf(colors.Bold+"%d."+colors.Reset+colors.Magenta+" %s\n", i+1, f+colors.Reset))
	}
	if len(files) == 0 {
		sb.WriteString(colors.Red + "No markdown files found.\n" + colors.Reset)
	}
	return sb.String()
}
