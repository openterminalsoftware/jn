package commands

import (
	"fmt"
	"jn/utils"
	"os"
	"path/filepath"
	"strings"
)

func Delete(file string) error {
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
			return fmt.Errorf("failed to get home directory: %w", err)
		}
		vaultPath = filepath.Join(homeDir, vaultPath[2:])
	} else if strings.HasPrefix(vaultPath, "$HOME/") {
		vaultPath = os.ExpandEnv(vaultPath)
	}

	var fullPath string
	if filepath.IsAbs(file) {
		fullPath = file
	} else {
		// If the file doesn't have a .md extension, add it
		if !strings.HasSuffix(file, ".md") {
			file = file + ".md"
		}
		fullPath = filepath.Join(vaultPath, file)
	}

	fullPath = os.ExpandEnv(fullPath)
	if strings.HasPrefix(fullPath, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}
		fullPath = filepath.Join(homeDir, fullPath[2:])
	}

	if err := os.Remove(fullPath); err != nil {
		return fmt.Errorf("failed to delete %s: %w", fullPath, err)
	}

	fmt.Printf("Deleted %s\n", fullPath)
	return nil
}
