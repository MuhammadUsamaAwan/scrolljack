package utils

import (
	"fmt"
	"os"

	"github.com/gen2brain/go-unarr"
)

func ExtractArchive(archivePath string, destinationPath string) error {
	a, err := unarr.NewArchive(archivePath)
	if err != nil {
		return fmt.Errorf("failed to open archive: %w", err)
	}
	defer a.Close()
	if err := os.MkdirAll(destinationPath, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}
	_, err = a.Extract(destinationPath)
	if err != nil {
		return fmt.Errorf("failed to extract archive: %w", err)
	}
	return nil
}
