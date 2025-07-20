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

func GetProfileFilesByProfileId(ctx context.Context, db *sql.DB, profileId string) ([]models.ProfileFile, error) {
	query := "SELECT id, profile_id, name, file_path FROM profile_files WHERE profile_id = $1"
	rows, err := db.QueryContext(ctx, query, profileId)
	if err != nil {
		return nil, fmt.Errorf("failed to query profile files: %w", err)
	}
	defer rows.Close()

	var profileFiles []models.ProfileFile
	for rows.Next() {
		var profileFile models.ProfileFile
		if err := rows.Scan(&profileFile.ID, &profileFile.ProfileID, &profileFile.Name, &profileFile.FilePath); err != nil {
			return nil, fmt.Errorf("failed to scan profile file row: %w", err)
		}
		profileFiles = append(profileFiles, profileFile)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred during row iteration: %w", err)
	}

	return profileFiles, nil
}
