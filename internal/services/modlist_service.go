package services

import (
	"context"
	"database/sql"
	"fmt"
	modlist "scrolljack/internal/types"
)

func InsertModlist(ctx context.Context, db *sql.DB, modlistId string, modlist *modlist.Modlist) error {
	_, err := db.ExecContext(
		ctx,
		`INSERT INTO modlists (id, name, author, description, game_type, image, readme, website, version, is_nsfw)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		modlistId,
		modlist.Name,
		modlist.Author,
		modlist.Description,
		modlist.GameType,
		modlist.Image,
		modlist.Readme,
		modlist.Website,
		modlist.Version,
		modlist.IsNSFW,
	)
	if err != nil {
		return fmt.Errorf("failed to insert modlist into database: %w", err)
	}
	return nil
}
