package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func ParseConfig(fileLocation string) map[string]interface{} {
	if len(fileLocation) >= 6 && fileLocation[:6] == "$HOME/" {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			fileLocation = homeDir + fileLocation[5:]
		}
	}

	data, err := ioutil.ReadFile(fileLocation)
	if err != nil {
		panic("Failed to read config file: " + err.Error())
	}

	// Attempt to fix common JSON mistakes, such as missing quotes around keys
	// For example: { vault: "~/.jn/vault" } -> { "vault": "~/.jn/vault" }
	trimmed := strings.TrimSpace(string(data))
	if len(trimmed) > 0 && trimmed[0] == '{' {
		// Replace unquoted keys with quoted keys
		// This is a naive fix and may not cover all cases
		fixed := ""
		inQuotes := false
		for i := 0; i < len(trimmed); i++ {
			c := trimmed[i]
			if c == '"' {
				inQuotes = !inQuotes
				fixed += string(c)
			} else if !inQuotes && ((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_') {
				// Start of a key
				start := i
				for i < len(trimmed) && ((trimmed[i] >= 'a' && trimmed[i] <= 'z') || (trimmed[i] >= 'A' && trimmed[i] <= 'Z') || trimmed[i] == '_' || (trimmed[i] >= '0' && trimmed[i] <= '9')) {
					i++
				}
				key := trimmed[start:i]
				fixed += "\"" + key + "\""
				// If next char is not ':', add it back
				if i < len(trimmed) && trimmed[i] != ':' {
					fixed += string(trimmed[i])
					i++
				}
				i--
			} else {
				fixed += string(c)
			}
		}
		data = []byte(fixed)
	}

	var config map[string]interface{}
	if err := json.Unmarshal(data, &config); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse config file: %v\n", err)
		fmt.Fprintf(os.Stderr, "Config file contents:\n%s\n", string(data))
		panic("Failed to parse config file: " + err.Error())
	}

	return config
}
