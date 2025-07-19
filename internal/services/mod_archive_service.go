package services

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"scrolljack/internal/db/models"
	modlist "scrolljack/internal/types"
	"scrolljack/internal/utils"
	"strings"
	"sync"

	"github.com/google/uuid"
)

func InsertModArchives(ctx context.Context, db *sql.DB, mods []models.Mod, m *modlist.Modlist) ([]models.ModArchive, error) {
	const chunkSize = 1000
	var directURLRegex = regexp.MustCompile(`(?m)^directURL=(.*)$`)

	archivesByHash := make(map[string]modlist.Archive)
	for _, archive := range m.Archives {
		archivesByHash[archive.Hash] = archive
	}

	archivesByName := make(map[string][]modlist.Archive)
	for _, archive := range m.Archives {
		if archive.State != nil && archive.State.Name != nil {
			name := utils.DerefStr(archive.State.Name)
			archivesByName[name] = append(archivesByName[name], archive)
		}
	}

	modArchivesChan := make(chan []models.ModArchive, len(mods))
	var wg sync.WaitGroup

	for _, mod := range mods {
		wg.Add(1)
		go func(mod models.Mod) {
			defer wg.Done()
			var modArchives []models.ModArchive
			var rawModArchives []modlist.Archive

			modPathPrefix := fmt.Sprintf("mods\\%s\\", mod.Name)

			for _, directive := range m.Directives {
				if strings.HasPrefix(directive.To, modPathPrefix) &&
					!strings.HasSuffix(directive.To, "meta.ini") {

					switch directive.Type {
					case modlist.FromArchiveType, modlist.PatchedFromArchiveType:
						if len(directive.ArchiveHashPath) > 0 {
							hash := directive.ArchiveHashPath[0]
							if archive, exists := archivesByHash[hash]; exists {
								rawModArchives = append(rawModArchives, archive)
							}
						}
					case modlist.CreateBSAType:
						if archives, exists := archivesByName[mod.Name]; exists {
							rawModArchives = append(rawModArchives, archives...)
						}
					}
				}
			}

			seen := make(map[string]struct{})
			for _, archive := range rawModArchives {
				if _, exists := seen[archive.Hash]; exists {
					continue
				}
				seen[archive.Hash] = struct{}{}

				directURL := ""
				if matches := directURLRegex.FindStringSubmatch(archive.Meta); len(matches) > 1 {
					directURL = matches[1]
				}

				typ := ""
				var (
					nexusGameName, version, description sql.NullString
					nexusModID, nexusFileID             sql.NullInt64
				)

				if archive.State != nil {
					typ = string(archive.State.Type)
					nexusGameName = utils.ToNullString(archive.State.GameName)
					version = utils.ToNullString(archive.State.Version)
					description = utils.ToNullString(archive.State.Description)
					nexusModID = utils.ToNullInt(archive.State.ModID)
					nexusFileID = utils.ToNullInt(archive.State.FileID)
				}

				info := models.ModArchive{
					ID:            uuid.New().String(),
					ModID:         mod.ID,
					Hash:          archive.Hash,
					Type:          typ,
					NexusGameName: nexusGameName,
					NexusModID:    nexusModID,
					NexusFileID:   nexusFileID,
					DirectURL:     utils.ToNullString(&directURL),
					Version:       version,
					Size:          utils.ToNullInt64(&archive.Size),
					Description:   description,
				}
				modArchives = append(modArchives, info)
			}
			modArchivesChan <- modArchives
		}(mod)
	}

	wg.Wait()
	close(modArchivesChan)

	var modArchivesToBeInserted []models.ModArchive
	for archives := range modArchivesChan {
		modArchivesToBeInserted = append(modArchivesToBeInserted, archives...)
	}

	if len(modArchivesToBeInserted) == 0 {
		return nil, nil
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction while inserting mod archives: %w", err)
	}
	defer tx.Rollback()

	for i := 0; i < len(modArchivesToBeInserted); i += chunkSize {
		chunkEnd := min(i+chunkSize, len(modArchivesToBeInserted))
		chunk := modArchivesToBeInserted[i:chunkEnd]

		var (
			valueStrings []string
			valueArgs    []any
		)
		for _, archive := range chunk {
			valueStrings = append(valueStrings, "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
			valueArgs = append(valueArgs,
				archive.ID,
				archive.ModID,
				archive.Hash,
				archive.Type,
				archive.NexusGameName,
				archive.NexusModID,
				archive.NexusFileID,
				archive.DirectURL,
				archive.Version,
				archive.Size,
				archive.Description,
			)
		}

		query := fmt.Sprintf(`
            INSERT INTO mod_archives (
                id, mod_id, hash, type, nexus_game_name, nexus_mod_id, nexus_file_id,
                direct_url, version, size, description
            ) VALUES %s`,
			strings.Join(valueStrings, ","),
		)

		if _, err := tx.ExecContext(ctx, query, valueArgs...); err != nil {
			return nil, fmt.Errorf("failed to insert mod archives in database: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("transaction commit failed while inserting mod archives: %w", err)
	}

	return modArchivesToBeInserted, nil
}
