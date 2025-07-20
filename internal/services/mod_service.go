package services

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"scrolljack/internal/db/dtos"
	"scrolljack/internal/db/models"
	modlist "scrolljack/internal/types"

	"github.com/google/uuid"
)

func InsertMods(ctx context.Context, db *sql.DB, profiles *[]models.Profile, modlist *modlist.Modlist, baseModlistPath string) ([]models.Mod, error) {
	const chunkSize = 1000
	var modsToBeInserted []models.Mod

	for _, profile := range *profiles {
		directive := findProfileModlistDirective(modlist, profile.Name)
		if directive == nil || directive.SourceDataID == nil || *directive.SourceDataID == "" {
			return nil, fmt.Errorf("modlist.txt not found for profile %s", profile.Name)
		}

		modlistFilePath := filepath.Join(baseModlistPath, *directive.SourceDataID)
		rawMods, err := readModlistFile(modlistFilePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read modlist.txt for profile %s: %w", profile.Name, err)
		}

		modOrder := 1
		for i, modName := range rawMods {
			isSeparator := strings.HasSuffix(modName, "_separator")
			modsToBeInserted = append(modsToBeInserted, models.Mod{
				ID:          uuid.New().String(),
				ProfileID:   profile.ID,
				Name:        strings.TrimSuffix(strings.TrimPrefix(modName, string(modName[0])), "_separator"),
				IsSeparator: isSeparator,
				Order:       i + 1,
				ModOrder:    0,
				IsActive:    strings.HasPrefix(modName, "+"),
			})
			if !isSeparator {
				modsToBeInserted[len(modsToBeInserted)-1].ModOrder = modOrder
				modOrder++
			}
		}
	}

	if len(modsToBeInserted) == 0 {
		return []models.Mod{}, nil
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction while inserting mods: %w", err)
	}
	defer tx.Rollback()

	for i := 0; i < len(modsToBeInserted); i += chunkSize {
		chunkEnd := min(i+chunkSize, len(modsToBeInserted))
		chunk := modsToBeInserted[i:chunkEnd]

		var (
			valueStrings []string
			valueArgs    []any
		)
		for _, mod := range chunk {
			valueStrings = append(valueStrings, "(?, ?, ?, ?, ?, ?, ?)")
			valueArgs = append(valueArgs,
				mod.ID,
				mod.ProfileID,
				mod.Name,
				mod.IsSeparator,
				mod.Order,
				mod.ModOrder,
				mod.IsActive,
			)
		}
		query := fmt.Sprintf(`
            INSERT INTO mods (id, profile_id, name, is_separator, "order", mod_order, is_active)
            VALUES %s`,
			strings.Join(valueStrings, ","),
		)

		if _, err := tx.ExecContext(ctx, query, valueArgs...); err != nil {
			return nil, fmt.Errorf("failed to insert mods in database: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("transaction commit failed while inserting mods: %w", err)
	}

	return modsToBeInserted, nil
}

func findProfileModlistDirective(modlist *modlist.Modlist, profileName string) *modlist.Directive {
	searchPath := fmt.Sprintf("profiles\\%s\\modlist.txt", profileName)
	for _, d := range modlist.Directives {
		if d.To == searchPath {
			return &d
		}
	}
	return nil
}

func readModlistFile(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append([]string{line}, lines...)
		}
	}
	return lines, scanner.Err()
}

func GetModsByProfileId(ctx context.Context, db *sql.DB, profileID string) ([]dtos.GroupedMod, error) {
	query := `
		SELECT id, profile_id, name, "order", is_separator
		FROM mods
		WHERE profile_id = ?
		ORDER BY "order" ASC
	`
	rows, err := db.Query(query, profileID)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	var groupedMods []dtos.GroupedMod
	var currentGroup *dtos.GroupedMod

	for rows.Next() {
		var mod dtos.Mod
		if err := rows.Scan(&mod.ID, &mod.ProfileID, &mod.Name, &mod.Order, &mod.IsSeparator); err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}

		if mod.IsSeparator {
			newGroup := dtos.GroupedMod{
				Separator: mod.Name,
				Mods:      make([]dtos.Mod, 0),
			}
			groupedMods = append(groupedMods, newGroup)
			currentGroup = &groupedMods[len(groupedMods)-1]
		} else {
			if currentGroup == nil {
				ungrouped := dtos.GroupedMod{
					Separator: "Ungrouped",
					Mods:      make([]dtos.Mod, 0),
				}
				groupedMods = append(groupedMods, ungrouped)
				currentGroup = &groupedMods[len(groupedMods)-1]
			}
			currentGroup.Mods = append(currentGroup.Mods, mod)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return groupedMods, nil
}
