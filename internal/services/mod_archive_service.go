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

	"github.com/google/uuid"
)

func InsertModArchives(ctx context.Context, db *sql.DB, mods []models.Mod, m *modlist.Modlist) ([]models.ModArchive, error) {
	const chunkSize = 1000
	var modArchivesToBeInserted []models.ModArchive
	var directURLRegex = regexp.MustCompile(`(?m)^directURL=(.*)$`)

	for _, mod := range mods {
		var rawModArchives []modlist.Archive

		for _, directive := range m.Directives {
			if strings.HasPrefix(directive.To, fmt.Sprintf("mods\\%s\\", mod.Name)) &&
				!strings.HasSuffix(directive.To, "meta.ini") {

				switch directive.Type {
				case "FromArchive", "PatchedFromArchive":
					if len(directive.ArchiveHashPath) > 0 {
						hash := directive.ArchiveHashPath[0]
						for _, archive := range m.Archives {
							if archive.Hash == hash {
								rawModArchives = append(rawModArchives, archive)
								break
							}
						}
					}
				case "CreateBSA":
					for _, archive := range m.Archives {
						if archive.State != nil && *archive.State.Name == mod.Name {
							rawModArchives = append(rawModArchives, archive)
						}
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

			info := models.ModArchive{
				ID:            uuid.New().String(),
				ModID:         mod.ID,
				Type:          string(archive.State.Type),
				NexusGameName: utils.ToNullString(archive.State.GameName),
				NexusModID:    utils.ToNullInt(archive.State.ModID),
				NexusFileID:   utils.ToNullInt(archive.State.FileID),
				DirectURL:     utils.ToNullString(&directURL),
				Version:       utils.ToNullString(archive.State.Version),
				Size:          utils.ToNullInt64(&archive.Size),
				Description:   utils.ToNullString(archive.State.Description),
			}
			modArchivesToBeInserted = append(modArchivesToBeInserted, info)
		}
	}

	if len(modArchivesToBeInserted) == 0 {
		return []models.ModArchive{}, nil
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
			valueStrings = append(valueStrings, "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
			valueArgs = append(valueArgs,
				archive.ID,
				archive.ModID,
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
            id, mod_id, type, nexus_game_name, nexus_mod_id, nexus_file_id,
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
