package services

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"scrolljack/internal/db/models"
	modlist "scrolljack/internal/types"
	"strings"

	"github.com/google/uuid"
)

func InsertProfileFiles(ctx context.Context, db *sql.DB, profiles *[]models.Profile, modlist *modlist.Modlist, baseModlistPath string) error {
	var profileFilesToBeInserted []models.ProfileFile

	for _, profile := range *profiles {
		for _, directive := range modlist.Directives {
			to := directive.To

			if !strings.HasPrefix(to, "profiles\\"+profile.Name+"\\") {
				continue
			}

			lowerTo := strings.ToLower(to)
			if strings.Contains(lowerTo, "backup") || strings.Contains(lowerTo, "cache") {
				continue
			}

			parts := strings.Split(to, "\\")
			name := parts[len(parts)-1]

			if directive.SourceDataID == nil || *directive.SourceDataID == "" {
				continue
			}

			profileFile := models.ProfileFile{
				ID:        uuid.New().String(),
				ProfileID: profile.ID,
				Name:      name,
				FilePath:  filepath.Join(baseModlistPath, *directive.SourceDataID),
			}

			profileFilesToBeInserted = append(profileFilesToBeInserted, profileFile)
		}
	}

	if len(profileFilesToBeInserted) == 0 {
		return nil
	}

	var (
		valueStrings []string
		valueArgs    []any
	)

	for _, pf := range profileFilesToBeInserted {
		valueStrings = append(valueStrings, "(?, ?, ?, ?)")
		valueArgs = append(valueArgs, pf.ID, pf.ProfileID, pf.Name, pf.FilePath)
	}

	query := fmt.Sprintf(`INSERT INTO profile_files (id, profile_id, name, file_path) VALUES %s`, strings.Join(valueStrings, ","))

	_, err := db.ExecContext(ctx, query, valueArgs...)
	if err != nil {
		return fmt.Errorf("failed to insert profile files into database: %w", err)
	}

	return nil
}
