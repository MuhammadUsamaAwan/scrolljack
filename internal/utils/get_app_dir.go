package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetAppDir() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user config directory: %w", err)
	}
	appDir := filepath.Join(configDir, "scrolljack")
	err = os.MkdirAll(appDir, 0755)
	if err != nil {
		return "", fmt.Errorf("failed to create app directory: %w", err)
	}
	return appDir, nil
}
