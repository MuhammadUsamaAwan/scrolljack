package services

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"scrolljack/internal/utils"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func DetectFomodOptions(ctx context.Context, db *sql.DB, modId string) (string, error) {
	result, err := runtime.OpenFileDialog(ctx, runtime.OpenDialogOptions{
		Title: "Select a mod archive (zip, rar, 7z)",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Mod Archive",
				Pattern:     "*.zip;*.rar;*.7z",
			},
		},
	})
	if err != nil {
		return "", fmt.Errorf("failed to open file dialog: %w", err)
	}
	if result == "" {
		return "", nil
	}

	appDir, err := utils.GetAppDir()
	if err != nil {
		return "", fmt.Errorf("failed to get app directory: %w", err)
	}
	tempDir := filepath.Join(appDir, "temp")
	if err := utils.ExtractArchive(result, tempDir); err != nil {
		return "", fmt.Errorf("failed to extract file: %w", err)
	}

	_, err = GetModFilesByModId(ctx, db, modId)
	if err != nil {
		return "", fmt.Errorf("failed to get mod files for mod ID %s: %w", modId, err)
	}

	fomodDir, moduleConfigPath, err := utils.FindFomodDirectory(tempDir)
	if err != nil {
		go cleanupTempDir(tempDir)
		return "", fmt.Errorf("failed to find FOMOD directory: %w", err)
	}

	if fomodDir == "" || moduleConfigPath == "" {
		go cleanupTempDir(tempDir)
		fmt.Println("No FOMOD directory or ModuleConfig.xml found in archive")
		return "No FOMOD configuration found in this mod archive", nil
	}

	go cleanupTempDir(tempDir)

	return "", nil
}

func cleanupTempDir(tempDir string) {
	if err := os.RemoveAll(tempDir); err != nil {
		log.Printf("Failed to delete temporary directory %s: %v", tempDir, err)
	}
}
