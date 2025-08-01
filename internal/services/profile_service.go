package services

import (
	"context"
	"database/sql"
	"fmt"
	"scrolljack/internal/db/models"
	modlist "scrolljack/internal/types"
	"strings"

	"github.com/google/uuid"
)

func InsertProfile(ctx context.Context, db *sql.DB, modlistId string, modlist *modlist.Modlist) ([]models.Profile, error) {
	var profilesToBeInserted []models.Profile
	for _, d := range modlist.Directives {
		if strings.HasPrefix(d.To, "profiles\\") && strings.HasSuffix(d.To, "\\modlist.txt") {
			parts := strings.Split(d.To, "\\")
			if len(parts) > 1 {
				newProfile := models.Profile{
					ID:        uuid.New().String(),
					ModlistID: modlistId,
					Name:      parts[1],
				}
				profilesToBeInserted = append(profilesToBeInserted, newProfile)
			}
		}
	}

	if len(profilesToBeInserted) == 0 {
		return nil, nil
	}

	valueStrings := make([]string, 0, len(profilesToBeInserted))
	valueArgs := make([]any, 0, len(profilesToBeInserted)*3)
	for i, profile := range profilesToBeInserted {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d)", i*3+1, i*3+2, i*3+3))
		valueArgs = append(valueArgs, profile.ID, profile.ModlistID, profile.Name)
	}

	query := fmt.Sprintf("INSERT INTO profiles (id, modlist_id, name) VALUES %s", strings.Join(valueStrings, ","))
	_, err := db.ExecContext(ctx, query, valueArgs...)
	if err != nil {
		return nil, fmt.Errorf("failed to insert profiles into database: %w", err)
	}

	return profilesToBeInserted, nil
}

func GetProfilesByModlistId(ctx context.Context, db *sql.DB, modlistId string) ([]models.Profile, error) {
	query := "SELECT id, modlist_id, name FROM profiles WHERE modlist_id = $1"
	rows, err := db.QueryContext(ctx, query, modlistId)
	if err != nil {
		return nil, fmt.Errorf("failed to query profiles: %w", err)
	}
	defer rows.Close()

	var profiles []models.Profile
	for rows.Next() {
		var profile models.Profile
		if err := rows.Scan(&profile.ID, &profile.ModlistID, &profile.Name); err != nil {
			return nil, fmt.Errorf("failed to scan profile row: %w", err)
		}
		profiles = append(profiles, profile)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred during row iteration: %w", err)
	}

	return profiles, nil
}
