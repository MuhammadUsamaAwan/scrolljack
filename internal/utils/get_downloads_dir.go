package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

func GetDownloadDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	var downloadDir string
	switch runtime.GOOS {
	case "windows":
		downloadDir = filepath.Join(home, "Downloads")
	case "darwin":
		downloadDir = filepath.Join(home, "Downloads")
	case "linux":
		downloadDir = filepath.Join(home, "Downloads")
	default:
		return "", fmt.Errorf("unsupported platform")
	}

	return downloadDir, nil
}
