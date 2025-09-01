package utils

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

type Config map[string]interface{}

type Entry struct {
	Title    string
	Content  string
	FileName string
}

func WriteEntry(config Config, entry Entry) error {
	vaultPathIface, ok := config["vault"]
	if !ok {
		return os.ErrNotExist
	}
	vaultPath, ok := vaultPathIface.(string)
	if !ok {
		return os.ErrNotExist
	}

	if strings.HasPrefix(vaultPath, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		vaultPath = filepath.Join(homeDir, vaultPath[2:])
	}

	if entry.FileName == "" {
		return os.ErrInvalid
	}

	entryFilePath := filepath.Join(vaultPath, entry.FileName)

	if _, err := os.Stat(vaultPath); os.IsNotExist(err) {
		if err := os.MkdirAll(vaultPath, 0700); err != nil {
			return err
		}
	}

	if _, err := os.Stat(entryFilePath); os.IsNotExist(err) {
		f, err := os.OpenFile(entryFilePath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
		if err != nil {
			return err
		}
		f.Close()
	} else if err == nil {
		return os.ErrExist
	} else {
		return err
	}

	entryData, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(entryFilePath, entryData, 0600); err != nil {
		return err
	}

	return nil
}
