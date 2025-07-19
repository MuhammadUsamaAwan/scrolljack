package main

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"scrolljack/internal/db"
	"scrolljack/internal/services"
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

func (a *App) ProcessWabbajackFile() {
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
		runtime.EventsEmit(a.ctx, "progress_update", fmt.Sprintf("❌ Failed to open file dialog: %v", err))
		return
	}
	if result == "" {
		return
	}

	modlistId := uuid.New().String()

	appDir, err := utils.GetAppDir()
	if err != nil {
		runtime.EventsEmit(a.ctx, "progress_update", fmt.Sprintf("❌ Failed to get app directory: %v", err))
		return
	}

	// Extract the Wabbajack file to the modlists directory
	runtime.EventsEmit(a.ctx, "progress_update", "📦 Extracting file...")
	start := time.Now()
	path := filepath.Join(appDir, "modlists", modlistId)
	if err := utils.ExtractArchive(result, path); err != nil {
		runtime.EventsEmit(a.ctx, "progress_update", fmt.Sprintf("❌ Failed to extract file: %v", err))
		return
	}
	runtime.EventsEmit(a.ctx, "progress_update", fmt.Sprintf("✅ Extraction completed in %s", utils.FormatDuration(time.Since(start))))

	// Read the modlist file
	start = time.Now()
	runtime.EventsEmit(a.ctx, "progress_update", "📖 Reading modlist file...")
	modlist, err := utils.LoadModlist(path)
	if err != nil {
		log.Fatal(err)
	}
	runtime.EventsEmit(a.ctx, "progress_update", fmt.Sprintf("✅ Modlist read in %s", utils.FormatDuration(time.Since(start))))

	// Save the modlist to the database
	start = time.Now()
	runtime.EventsEmit(a.ctx, "progress_update", "💾 Saving modlist to database...")
	if err := services.InsertModlist(modlistId, modlist); err != nil {
		runtime.EventsEmit(a.ctx, "progress_update", fmt.Sprintf("Failed to save modlist: %v", err))
		return
	}
	runtime.EventsEmit(a.ctx, "progress_update", fmt.Sprintf("✅ Modlist saved in %s", utils.FormatDuration(time.Since(start))))

	// Save the profiles to the database
	start = time.Now()
	runtime.EventsEmit(a.ctx, "progress_update", "📂 Saving profiles to database...")
	profiles, err := services.InsertProfile(modlistId, modlist)
	if err != nil {
		runtime.EventsEmit(a.ctx, "progress_update", fmt.Sprintf("❌ Failed to save profiles: %v", err))
		return
	}
	runtime.EventsEmit(a.ctx, "progress_update", fmt.Sprintf("✅ %d profiles saved in %s", len(profiles), utils.FormatDuration(time.Since(start))))

	// Save the profile files to the database
	start = time.Now()
	runtime.EventsEmit(a.ctx, "progress_update", "📄 Saving profile files to database...")
	if err := services.InsertProfileFiles(&profiles, modlist, path); err != nil {
		runtime.EventsEmit(a.ctx, "progress_update", fmt.Sprintf("❌ Failed to save profile files: %v", err))
		return
	}
	runtime.EventsEmit(a.ctx, "progress_update", fmt.Sprintf("✅ Profile files saved in %s", utils.FormatDuration(time.Since(start))))

}
