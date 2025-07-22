package services

import (
	"bufio"
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"scrolljack/internal/utils"
	"unicode/utf8"

	"github.com/OctopusDeploy/go-octodiff/pkg/octodiff"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func BinaryPatch(ctx context.Context, PatchFilePath string, name string) (string, string, error) {
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
		return "", "", fmt.Errorf("failed to open file dialog: %w", err)
	}
	if result == "" {
		return "", "", nil
	}

	srcFile, err := os.Open(result)
	if err != nil {
		return "", "", fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcFile.Close()

	patchFile, err := os.Open(PatchFilePath)
	if err != nil {
		return "", "", fmt.Errorf("failed to open patch file: %w", err)
	}
	defer patchFile.Close()

	downloadsDir, err := utils.GetDownloadDir()
	if err != nil {
		return "", "", fmt.Errorf("failed to get downloads directory: %w", err)
	}

	dstPath := filepath.Join(downloadsDir, name)
	dstFile, err := os.Create(dstPath)
	if err != nil {
		return "", "", fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dstFile.Close()

	deltaReader := octodiff.NewBinaryDeltaReader(bufio.NewReader(patchFile))
	dstWriter := bufio.NewWriter(dstFile)

	err = octodiff.ApplyDelta(srcFile, deltaReader, dstWriter)
	if err != nil {
		return "", "", fmt.Errorf("failed to apply octodiff patch: %w", err)
	}

	if err := dstWriter.Flush(); err != nil {
		return "", "", fmt.Errorf("failed to write patched file: %w", err)
	}

	srcBase64, err := readFileBase64(result)
	if err != nil {
		return "", "", fmt.Errorf("error reading source file base64: %w", err)
	}

	dstBase64, err := readFileBase64(dstPath)
	if err != nil {
		return "", "", fmt.Errorf("error reading patched file base64: %w", err)
	}

	return srcBase64, dstBase64, nil
}

const maxFileSize int64 = 10 * 1024 * 1024 // 10 MB

func readFileBase64(path string) (string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return "", err
	}

	if info.Size() > maxFileSize {
		return "", nil
	}

	bytes, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	if !utf8.Valid(bytes) {
		return "", nil
	}

	encoded := base64.StdEncoding.EncodeToString(bytes)
	return encoded, nil
}
