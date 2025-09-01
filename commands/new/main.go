package new

import (
	"fmt"
	"jn/colors"
	"jn/utils"
	"os"
)

func New() {
	if _, err := os.Stat(os.ExpandEnv("$HOME/.jn/config.json")); os.IsNotExist(err) {
		fmt.Println("Creating ~/.jn/config.json")
		configPath := os.ExpandEnv("$HOME/.jn/config.json")
		if err := os.MkdirAll(os.ExpandEnv("$HOME/.jn"), 0755); err != nil {
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

	if _, err := os.Stat(os.ExpandEnv("$HOME/.jn/config.json")); err == nil {
		fmt.Println(colors.DarkGray + "Using ~/.jn/config.json" + colors.Reset)
		config := utils.ParseConfig(os.ExpandEnv("$HOME/.jn/config.json"))
		makeEntry(config)
	}
}

func makeEntry(config utils.Config) {
	fmt.Println("Creating new entry")
	entry := utils.Entry{
		Title:    "New Entry",
		Content:  "This is a new entry",
		FileName: os.Args[2],
	}
	if err := utils.WriteEntry(config, entry); err != nil {
		fmt.Printf("Failed to write entry: %v\n", err)
		return
	}
	fmt.Println("Entry created successfully")
}

func main() {
	New()
}
