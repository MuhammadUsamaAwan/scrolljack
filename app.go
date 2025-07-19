package main

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"scrolljack/internal/db"
	"scrolljack/internal/utils"
	"time"

	"github.com/google/uuid"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	db.Connect()
}

func (a *App) shutdown(ctx context.Context) {
	if db.DB != nil {
		db.DB.Close()
	}
}

func (a *App) ProcessWabbajackFile() error {
	result, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select a file",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Wabbajack File",
				Pattern:     "*.wabbajack",
			},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to open file dialog: %w", err)
	}

	modlistId := uuid.New().String()

	appDir, err := utils.GetAppDir()
	if err != nil {
		return err
	}

	// Extract the Wabbajack file to the modlists directory
	runtime.EventsEmit(a.ctx, "progress_update", "📦 Extracting file...")
	start := time.Now()
	path := filepath.Join(appDir, "modlists", modlistId)
	log.Printf("Extracting Wabbajack file to %s", path)
	if err := utils.ExtractArchive(result, path); err != nil {
		return err
	}
	runtime.EventsEmit(a.ctx, "progress_update", fmt.Sprintf("✅ Extraction completed in %s", time.Since(start).String()))

	return nil
}
