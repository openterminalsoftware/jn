package commands

import (
	"bufio"
	"encoding/json"
	"fmt"
	"jn/colors"
	"os"
	"path/filepath"
	"strings"
)

type ConfigData struct {
	Vault string `json:"vault"`
}

func Config() {
	scanner := bufio.NewScanner(os.Stdin)
	defaultVault := os.ExpandEnv("$HOME/.jn/vault")

	fmt.Print(colors.Bold + "Enter vault path" + colors.Reset)
	fmt.Printf(" (default: %s): ", defaultVault)

	scanner.Scan()
	input := scanner.Text()

	vaultPath := strings.TrimSpace(input)
	if vaultPath == "" {
		vaultPath = defaultVault
	}

	config := ConfigData{
		Vault: vaultPath,
	}

	configPath := os.ExpandEnv("$HOME/.jn/config.json")
	configDir := filepath.Dir(configPath)

	if err := os.MkdirAll(configDir, 0755); err != nil {
		fmt.Printf(colors.Red+"Failed to create config directory: %v\n"+colors.Reset, err)
		return
	}

	file, err := os.Create(configPath)
	if err != nil {
		fmt.Printf(colors.Red+"Failed to create config file: %v\n"+colors.Reset, err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(config); err != nil {
		fmt.Printf(colors.Red+"Failed to write to config file: %v\n"+colors.Reset, err)
		return
	}

	fmt.Println("\n" + colors.Green + "Configuration saved to " + colors.Bold + configPath + colors.Reset)
}
