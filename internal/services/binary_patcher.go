package services

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"scrolljack/internal/utils"

	"github.com/OctopusDeploy/go-octodiff/pkg/octodiff"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func BinaryPatch(ctx context.Context, PatchFilePath string, name string) error {
	// Ask user to select source file
	result, err := runtime.OpenFileDialog(ctx, runtime.OpenDialogOptions{
		Title: "Select a mod file to patch",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Mod File",
				Pattern:     "*",
			},
		},
	})
	if err != nil {
		log.Printf("Error opening file dialog: %v", err)
		return fmt.Errorf("failed to open file dialog: %w", err)
	}
	if result == "" {
		return nil
	}

	srcFile, err := os.Open(result)
	if err != nil {
		log.Printf("Error opening source file: %v", err)
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcFile.Close()

	patchFile, err := os.Open(PatchFilePath)
	if err != nil {
		log.Printf("Error opening patch file: %v", err)
		return fmt.Errorf("failed to open patch file: %w", err)
	}
	defer patchFile.Close()

	log.Printf("Applying octodiff patch from %s to %s", result, PatchFilePath)

	downloadsDir, err := utils.GetDownloadDir()
	if err != nil {
		log.Printf("Error getting downloads directory: %v", err)
		return fmt.Errorf("failed to get downloads directory: %w", err)
	}

	dstPath := filepath.Join(downloadsDir, name)
	dstFile, err := os.Create(dstPath)
	if err != nil {
		log.Printf("Error creating destination file: %v", err)
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dstFile.Close()

	deltaReader := octodiff.NewBinaryDeltaReader(bufio.NewReader(patchFile))
	dstWriter := bufio.NewWriter(dstFile)

	// Apply octodiff patch
	err = octodiff.ApplyDelta(srcFile, deltaReader, dstWriter)
	if err != nil {
		log.Printf("Error applying octodiff patch: %v", err)
		return fmt.Errorf("failed to apply octodiff patch: %w", err)
	}

	if err := dstWriter.Flush(); err != nil {
		log.Printf("Error flushing output file: %v", err)
		return fmt.Errorf("failed to write patched file: %w", err)
	}

	return nil
}
