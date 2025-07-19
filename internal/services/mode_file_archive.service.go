package services

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"sync"

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
	hashToArchiveID := make(map[string]string, len(modArchives))
	for _, archive := range modArchives {
		hashToArchiveID[archive.Hash] = archive.ID
	}

	modFileMap := make(map[string]map[string]string, len(modFiles))
	for _, mf := range modFiles {
		if modFileMap[mf.ModID] == nil {
			modFileMap[mf.ModID] = make(map[string]string)
		}
		modFileMap[mf.ModID][mf.Path] = mf.ID
	}

	const chunkSize = 1000
	var linksMu sync.Mutex
	var linksToInsert []models.ModFileArchive

	var wg sync.WaitGroup
	wg.Add(len(mods))

	for _, mod := range mods {
		mod := mod // capture range variable
		go func() {
			defer wg.Done()
			var links []models.ModFileArchive
			for _, directive := range modlist.Directives {
				if !strings.HasPrefix(directive.To, fmt.Sprintf("mods\\%s\\", mod.Name)) ||
					strings.HasSuffix(directive.To, "meta.ini") {
					continue
				}
				if len(directive.ArchiveHashPath) == 0 {
					continue
				}

				archiveID, ok := hashToArchiveID[directive.ArchiveHashPath[0]]
				if !ok {
					continue
				}
				relativePath := strings.Replace(directive.To, fmt.Sprintf("mods\\%s\\", mod.Name), "", 1)

				if filesForMod, ok := modFileMap[mod.ID]; ok {
					if fileID, ok := filesForMod[relativePath]; ok {
						links = append(links, models.ModFileArchive{
							ModlistId:    modlistID,
							ModFileId:    fileID,
							ModArchiveId: archiveID,
						})
					}
				}
			}
			if len(links) > 0 {
				linksMu.Lock()
				linksToInsert = append(linksToInsert, links...)
				linksMu.Unlock()
			}
		}()
	}
	wg.Wait()

	if len(linksToInsert) == 0 {
		return nil
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction for mod file archive links: %w", err)
	}
	defer tx.Rollback()

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
