package services

import (
	"context"
	"fmt"
	"scrolljack/internal/db"
	"scrolljack/internal/db/models"
	modlist "scrolljack/internal/types"
	"strings"

	"github.com/google/uuid"
)

func InsertProfile(modlistId string, modlist *modlist.Modlist) ([]models.Profile, error) {
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
	_, err := db.DB.ExecContext(context.Background(), query, valueArgs...)
	if err != nil {
		return nil, fmt.Errorf("failed to insert profiles into database: %w", err)
	}

	return profilesToBeInserted, nil
}
