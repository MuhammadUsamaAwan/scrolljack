package services

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"scrolljack/internal/db/dtos"
	"scrolljack/internal/db/models"
	modlist "scrolljack/internal/types"
	"scrolljack/internal/utils"

	"github.com/google/uuid"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func InsertModFiles(ctx context.Context, db *sql.DB, mods []models.Mod, m *modlist.Modlist, baseModlistPath string) ([]models.ModFile, error) {
	const chunkSize = 1000

	directivesByMod := make(map[string][]modlist.Directive)
	for _, directive := range m.Directives {
		if strings.HasPrefix(directive.To, "mods\\") && !strings.HasSuffix(directive.To, "meta.ini") {
			parts := strings.Split(directive.To, "\\")
			if len(parts) >= 2 {
				modName := parts[1]
				directivesByMod[modName] = append(directivesByMod[modName], directive)
			}
		}
	}

	modFilesChan := make(chan []models.ModFile, len(mods))
	var wg sync.WaitGroup

	for _, mod := range mods {
		wg.Add(1)
		go func(mod models.Mod) {
			defer wg.Done()

			rawModFiles, exists := directivesByMod[mod.Name]
			if !exists {
				modFilesChan <- []models.ModFile{}
				return
			}

			modFiles := make([]models.ModFile, 0, len(rawModFiles))

			modPathPrefix := fmt.Sprintf("mods\\%s\\", mod.Name)

			for _, mf := range rawModFiles {
				var sourceFilePath sql.NullString
				var patchFilePath sql.NullString

				if mf.SourceDataID != nil && *mf.SourceDataID != "" {
					fullPath := filepath.Join(baseModlistPath, *mf.SourceDataID)
					sourceFilePath = utils.ToNullString(&fullPath)
				}

				if mf.PatchID != nil && *mf.PatchID != "" {
					fullPath := filepath.Join(baseModlistPath, *mf.PatchID)
					patchFilePath = utils.ToNullString(&fullPath)
				}

				fileStatePtrs := make([]*modlist.FileState, 0, len(mf.FileStates))
				for i := range mf.FileStates {
					fileStatePtrs = append(fileStatePtrs, &mf.FileStates[i])
				}
				bsaFilesStr := strings.Join(extractPaths(fileStatePtrs), ";")

				relativePath := strings.TrimPrefix(mf.To, modPathPrefix)

				modFile := models.ModFile{
					ID:             uuid.New().String(),
					ModID:          mod.ID,
					Hash:           mf.Hash,
					Type:           string(mf.Type),
					Path:           relativePath,
					SourceFilePath: sourceFilePath,
					PatchFilePath:  patchFilePath,
					BsaFiles:       utils.ToNullString(&bsaFilesStr),
				}
				modFiles = append(modFiles, modFile)
			}
			modFilesChan <- modFiles
		}(mod)
	}

	wg.Wait()
	close(modFilesChan)

	var modFilesToBeInserted []models.ModFile
	for files := range modFilesChan {
		modFilesToBeInserted = append(modFilesToBeInserted, files...)
	}

	if len(modFilesToBeInserted) == 0 {
		return []models.ModFile{}, nil
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction while inserting mod files: %w", err)
	}
	defer tx.Rollback()

	for i := 0; i < len(modFilesToBeInserted); i += chunkSize {
		chunkEnd := min(i+chunkSize, len(modFilesToBeInserted))
		chunk := modFilesToBeInserted[i:chunkEnd]

		var (
			valueStrings []string
			valueArgs    []any
		)
		for _, file := range chunk {
			valueStrings = append(valueStrings, "(?, ?, ?, ?, ?, ?, ?, ?)")
			valueArgs = append(valueArgs,
				file.ID,
				file.ModID,
				file.Hash,
				file.Type,
				file.Path,
				file.SourceFilePath,
				file.PatchFilePath,
				file.BsaFiles,
			)
		}

		query := fmt.Sprintf(`
        INSERT INTO mod_files (
            id, mod_id, hash, type, path, source_file_path, patch_file_path, bsa_files
        ) VALUES %s`,
			strings.Join(valueStrings, ","),
		)

		if _, err := tx.ExecContext(ctx, query, valueArgs...); err != nil {
			return nil, fmt.Errorf("failed to insert mod files in database: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("transaction commit failed while inserting mod files: %w", err)
	}

	return modFilesToBeInserted, nil
}

func extractPaths(states []*modlist.FileState) []string {
	paths := make([]string, 0, len(states))
	for _, fs := range states {
		if fs != nil {
			paths = append(paths, fs.Path)
		}
	}
	return paths
}

func GetModFilesByModId(ctx context.Context, db *sql.DB, modID string) ([]dtos.ModFileDTO, error) {
	query := `
		SELECT id, hash, type, path, source_file_path, patch_file_path, bsa_files
		FROM mod_files
		WHERE mod_id = $1
	`

	rows, err := db.QueryContext(ctx, query, modID)
	if err != nil {
		return nil, fmt.Errorf("failed to query mod files for mod ID %s: %w", modID, err)
	}
	defer rows.Close()

	var modFiles []dtos.ModFileDTO
	for rows.Next() {
		var file dtos.ModFileDTO
		if err := rows.Scan(&file.ID, &file.Hash, &file.Type, &file.Path,
			&file.SourceFilePath, &file.PatchFilePath, &file.BsaFiles); err != nil {
			return nil, fmt.Errorf("failed to scan mod file row: %w", err)
		}
		modFiles = append(modFiles, file)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred while iterating over mod files: %w", err)
	}

	return modFiles, nil
}

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
	log.Printf("[DEBUG] Extracting archive to: %s\n", tempDir)
	if err := utils.ExtractArchive(result, tempDir); err != nil {
		return "", fmt.Errorf("failed to extract file: %w", err)
	}

	fomodDir := filepath.Join(tempDir, "fomod")
	if _, err := os.Stat(fomodDir); os.IsNotExist(err) {
		return "", fmt.Errorf("FOMOD directory not found: %s", fomodDir)
	}
	moduleConfigPath := filepath.Join(fomodDir, "ModuleConfig.xml")
	if _, err := os.Stat(moduleConfigPath); os.IsNotExist(err) {
		return "", fmt.Errorf("no ModuleConfig.xml found: %s", moduleConfigPath)
	}

	_, err = GetModFilesByModId(ctx, db, modId)
	if err != nil {
		return "", fmt.Errorf("failed to get mod files for mod ID %s: %w", modId, err)
	}

	// TODO: Implement FOMOD option detection logic here

	if err := os.RemoveAll(tempDir); err != nil {
		return "", fmt.Errorf("failed to delete temporary directory: %w", err)
	}

	return "", nil
}
