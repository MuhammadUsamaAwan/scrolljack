package main

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"scrolljack/internal/db"
	"scrolljack/internal/db/dtos"
	"scrolljack/internal/db/models"
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
	globalStart := time.Now()
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
	if err := services.InsertModlist(a.ctx, db.DB, modlistId, modlist); err != nil {
		runtime.EventsEmit(a.ctx, "progress_update", fmt.Sprintf("Failed to save modlist: %v", err))
		return
	}
	runtime.EventsEmit(a.ctx, "progress_update", fmt.Sprintf("✅ Modlist saved in %s", utils.FormatDuration(time.Since(start))))

	// Save the profiles to the database
	start = time.Now()
	runtime.EventsEmit(a.ctx, "progress_update", "📂 Saving profiles to database...")
	profiles, err := services.InsertProfile(a.ctx, db.DB, modlistId, modlist)
	if err != nil {
		runtime.EventsEmit(a.ctx, "progress_update", fmt.Sprintf("❌ Failed to save profiles: %v", err))
		return
	}
	runtime.EventsEmit(a.ctx, "progress_update", fmt.Sprintf("✅ %d profiles saved in %s", len(profiles), utils.FormatDuration(time.Since(start))))

	// Save the profile files to the database
	start = time.Now()
	runtime.EventsEmit(a.ctx, "progress_update", "📄 Saving profile files to database...")
	if err := services.InsertProfileFiles(a.ctx, db.DB, &profiles, modlist, path); err != nil {
		runtime.EventsEmit(a.ctx, "progress_update", fmt.Sprintf("❌ Failed to save profile files: %v", err))
		return
	}
	runtime.EventsEmit(a.ctx, "progress_update", fmt.Sprintf("✅ Profile files saved in %s", utils.FormatDuration(time.Since(start))))

	// Save the mods to the database
	start = time.Now()
	runtime.EventsEmit(a.ctx, "progress_update", "🔧 Saving mods to database...")
	mods, err := services.InsertMods(a.ctx, db.DB, &profiles, modlist, path)
	if err != nil {
		runtime.EventsEmit(a.ctx, "progress_update", fmt.Sprintf("❌ Failed to save mods: %v", err))
		return
	}
	runtime.EventsEmit(a.ctx, "progress_update", fmt.Sprintf("✅ Mods saved in %s", utils.FormatDuration(time.Since(start))))

	// Save the mod archives to the database
	start = time.Now()
	runtime.EventsEmit(a.ctx, "progress_update", "📦 Saving mod archives to database...")
	archives, err := services.InsertModArchives(a.ctx, db.DB, mods, modlist)
	if err != nil {
		runtime.EventsEmit(a.ctx, "progress_update", fmt.Sprintf("❌ Failed to save mod archives: %v", err))
		return
	}
	runtime.EventsEmit(a.ctx, "progress_update", fmt.Sprintf("✅ Mod archives saved in %s", utils.FormatDuration(time.Since(start))))

	// Save the mod files to the database
	start = time.Now()
	runtime.EventsEmit(a.ctx, "progress_update", "📂 Saving mod files to database...")
	files, err := services.InsertModFiles(a.ctx, db.DB, mods, modlist, path)
	if err != nil {
		runtime.EventsEmit(a.ctx, "progress_update", fmt.Sprintf("❌ Failed to save mod files: %v", err))
		return
	}
	runtime.EventsEmit(a.ctx, "progress_update", fmt.Sprintf("✅ Mod files saved in %s", utils.FormatDuration(time.Since(start))))

	// Save the mod file archive links to the database
	start = time.Now()
	runtime.EventsEmit(a.ctx, "progress_update", "🔗 Saving mod file archive links...")
	if err := services.InsertModFileArchiveLinks(a.ctx, db.DB, modlistId, mods, files, archives, modlist); err != nil {
		runtime.EventsEmit(a.ctx, "progress_update", fmt.Sprintf("❌ Failed to save mod file archive links: %v", err))
		return
	}
	runtime.EventsEmit(a.ctx, "progress_update", fmt.Sprintf("✅ Mod file archive links saved in %s", utils.FormatDuration(time.Since(start))))

	runtime.EventsEmit(a.ctx, "progress_update", fmt.Sprintf("🎉 Modlist import completed in %s", utils.FormatDuration(time.Since(globalStart))))

}

func (a *App) GetModlists() ([]*dtos.ModlistDTO, error) {
	modlists, err := services.GetModlists(a.ctx, db.DB)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve modlists: %w", err)
	}
	return modlists, nil
}

func (a *App) GetModlistImageBase64(modlistId string, image string) (string, error) {
	base64Image, err := services.GetModlistImageBase64(modlistId, image)
	if err != nil {
		return "", fmt.Errorf("failed to get modlist image: %w", err)
	}
	return base64Image, nil
}

func (a *App) GetModlistById(modlistId string) (*dtos.ModlistDTO, error) {
	modlist, err := services.GetModlistById(a.ctx, db.DB, modlistId)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve modlist by ID: %w", err)
	}
	return modlist, nil
}

func (a *App) DeleteModlist(modlistId string) error {
	if err := services.DeleteModlist(a.ctx, db.DB, modlistId); err != nil {
		return fmt.Errorf("failed to delete modlist: %w", err)
	}
	return nil
}

func (a *App) GetProfilesByModlistId(modlistId string) ([]models.Profile, error) {
	profiles, err := services.GetProfilesByModlistId(a.ctx, db.DB, modlistId)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve profiles by modlist ID: %w", err)
	}
	return profiles, nil
}

func (a *App) GetProfileFilesByProfileId(profileId string) ([]models.ProfileFile, error) {
	profileFiles, err := services.GetProfileFilesByProfileId(a.ctx, db.DB, profileId)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve profile files by profile ID: %w", err)
	}
	return profileFiles, nil
}

func (a *App) DownloadFile(path string, name string) error {
	downloadsDir, err := utils.GetDownloadDir()
	if err != nil {
		return fmt.Errorf("failed to get downloads directory: %w", err)
	}
	dstPath := filepath.Join(downloadsDir, name)
	err = utils.CopyFile(path, dstPath)
	if err != nil {
		return fmt.Errorf("failed to copy profile file to downloads directory: %w", err)
	}
	return nil
}

func (a *App) GetModsByProfileId(profileId string) ([]dtos.GroupedModDTO, error) {
	groupedMods, err := services.GetModsByProfileId(a.ctx, db.DB, profileId)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve mods by profile ID: %w", err)
	}
	return groupedMods, nil
}

func (a *App) GetModArchivesByModId(modId string) ([]dtos.ModArchiveDTO, error) {
	archives, err := services.GetModArchivesByModId(a.ctx, db.DB, modId)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve mod archives by mod ID: %w", err)
	}
	return archives, nil
}

func (a *App) GetModFilesByModId(modId string) ([]dtos.ModFileDTO, error) {
	modFiles, err := services.GetModFilesByModId(a.ctx, db.DB, modId)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve mod files by mod ID: %w", err)
	}
	return modFiles, nil
}

func (a *App) DetectFomodOptions(modId string) (string, error) {
	fomodOptions, err := services.DetectFomodOptions(a.ctx, db.DB, modId)
	if err != nil {
		return "", fmt.Errorf("failed to get file differences: %w", err)
	}
	return fomodOptions, nil
}
