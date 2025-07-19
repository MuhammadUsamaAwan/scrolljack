package services

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"scrolljack/internal/db/models"
	modlist "scrolljack/internal/types"
)

func InsertModFileArchiveLinks(
	ctx context.Context,
	db *sql.DB,
	modlistID string,
	mods []models.Mod,
	modFiles []models.ModFile,
	modArchives []models.ModArchive,
	modlist *modlist.Modlist,
) error {
	hashToArchiveID := make(map[string]string)
	for _, archive := range modArchives {
		hashToArchiveID[archive.Hash] = archive.ID
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction for mod file archive links: %w", err)
	}
	defer tx.Rollback()

	const chunkSize = 1000
	var linksToInsert []models.ModFileArchive

	for _, mod := range mods {
		for _, directive := range modlist.Directives {
			if !strings.HasPrefix(directive.To, fmt.Sprintf("mods\\%s\\", mod.Name)) ||
				strings.HasSuffix(directive.To, "meta.ini") {
				continue
			}

			if len(directive.ArchiveHashPath) == 0 {
				continue
			}

			archiveHash := directive.ArchiveHashPath[0]
			archiveID, ok := hashToArchiveID[archiveHash]
			if !ok {
				continue
			}

			relativePath := strings.Replace(directive.To, fmt.Sprintf("mods\\%s\\", mod.Name), "", 1)

			for _, modFile := range modFiles {
				if modFile.ModID == mod.ID && modFile.Path == relativePath {
					linksToInsert = append(linksToInsert, models.ModFileArchive{
						ModlistId:    modlistID,
						ModFileId:    modFile.ID,
						ModArchiveId: archiveID,
					})
					break
				}
			}
		}
	}

	for i := 0; i < len(linksToInsert); i += chunkSize {
		chunkEnd := min(i+chunkSize, len(linksToInsert))
		chunk := linksToInsert[i:chunkEnd]

		var (
			valueStrings []string
			valueArgs    []any
		)
		for _, link := range chunk {
			valueStrings = append(valueStrings, "(?, ?, ?)")
			valueArgs = append(valueArgs, link.ModlistId, link.ModFileId, link.ModArchiveId)
		}

		query := fmt.Sprintf(`
            INSERT INTO mod_file_archives (modlist_id, mod_file_id, mod_archive_id)
            VALUES %s
        `, strings.Join(valueStrings, ","))

		if _, err := tx.ExecContext(ctx, query, valueArgs...); err != nil {
			return fmt.Errorf("failed to insert mod file archive links: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("transaction commit failed for mod file archive links: %w", err)
	}

	return nil
}
