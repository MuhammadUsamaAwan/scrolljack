package services

import (
	"context"
	"fmt"
	"scrolljack/internal/db"
	modlist "scrolljack/internal/types"
)

func InsertModlist(modlistId string, modlist *modlist.Modlist) error {
	_, err := db.DB.ExecContext(
		context.Background(),
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
