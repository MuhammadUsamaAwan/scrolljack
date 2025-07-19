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
	m *modlist.Modlist,
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

	directivesByMod := make(map[string][]modlist.Directive)
	for _, directive := range m.Directives {
		if len(directive.ArchiveHashPath) == 0 {
			continue
		}

		if strings.HasPrefix(directive.To, "mods\\") && !strings.HasSuffix(directive.To, "meta.ini") {
			parts := strings.Split(directive.To, "\\")
			if len(parts) >= 2 {
				modName := parts[1]
				directivesByMod[modName] = append(directivesByMod[modName], directive)
			}
		}
	}

	const chunkSize = 1000
	var linksMu sync.Mutex
	var linksToInsert []models.ModFileArchive

	var wg sync.WaitGroup

	for _, mod := range mods {
		relevantDirectives, hasDirectives := directivesByMod[mod.Name]
		if !hasDirectives {
			continue
		}

		wg.Add(1)
		go func(mod models.Mod, directives []modlist.Directive) {
			defer wg.Done()

			links := make([]models.ModFileArchive, 0, len(directives))

			modPathPrefix := fmt.Sprintf("mods\\%s\\", mod.Name)
			filesForMod, hasFiles := modFileMap[mod.ID]
			if !hasFiles {
				return
			}

			for _, directive := range directives {
				archiveID, archiveExists := hashToArchiveID[directive.ArchiveHashPath[0]]
				if !archiveExists {
					continue
				}

				relativePath := strings.TrimPrefix(directive.To, modPathPrefix)

				if fileID, fileExists := filesForMod[relativePath]; fileExists {
					links = append(links, models.ModFileArchive{
						ModlistId:    modlistID,
						ModFileId:    fileID,
						ModArchiveId: archiveID,
					})
				}
			}

			if len(links) > 0 {
				linksMu.Lock()
				linksToInsert = append(linksToInsert, links...)
				linksMu.Unlock()
			}
		}(mod, relevantDirectives)
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
