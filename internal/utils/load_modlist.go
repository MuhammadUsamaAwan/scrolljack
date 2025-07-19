package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	modlist "scrolljack/internal/types"
)

func LoadModlist(baseModlistPath string) (*modlist.Modlist, error) {
	path := filepath.Join(baseModlistPath, "modlist")
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open modlist file: %w", err)
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read modlist file: %w", err)
	}

	var ml modlist.Modlist
	if err := json.Unmarshal(bytes, &ml); err != nil {
		return nil, fmt.Errorf("failed to unmarshal modlist: %w", err)
	}

	return &ml, nil
}
