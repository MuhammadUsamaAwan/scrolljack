package services

import (
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"scrolljack/internal/db/dtos"
	modlist "scrolljack/internal/types"
	"scrolljack/internal/utils"
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

func GetModlists(ctx context.Context, db *sql.DB) ([]*dtos.ModlistDTO, error) {
	rows, err := db.QueryContext(ctx, `SELECT id, name, author, description, image, game_type, version, is_nsfw, website, readme, created_at FROM modlists`)
	if err != nil {
		return nil, fmt.Errorf("failed to query modlists: %w", err)
	}
	defer rows.Close()

	var modlists []*dtos.ModlistDTO
	for rows.Next() {
		var m dtos.ModlistDTO
		if err := rows.Scan(&m.ID, &m.Name, &m.Author, &m.Description, &m.Image, &m.GameType, &m.Version, &m.IsNSFW, &m.Website, &m.Readme, &m.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan modlist row: %w", err)
		}
		modlists = append(modlists, &m)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error encountered during row iteration: %w", err)
	}

	return modlists, nil
}

func GetModlistById(ctx context.Context, db *sql.DB, modlistId string) (*dtos.ModlistDTO, error) {
	row := db.QueryRowContext(ctx, `SELECT id, name, author, description, image, game_type, version, is_nsfw, website, readme, created_at FROM modlists WHERE id = ?`, modlistId)

	var m dtos.ModlistDTO
	if err := row.Scan(&m.ID, &m.Name, &m.Author, &m.Description, &m.Image, &m.GameType, &m.Version, &m.IsNSFW, &m.Website, &m.Readme, &m.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to scan modlist: %w", err)
	}

	return &m, nil
}

func GetModlistImageBase64(modlistId string, image string) (string, error) {
	appDir, err := utils.GetAppDir()
	if err != nil {
		return "", fmt.Errorf("failed to get app directory: %w", err)
	}
	path := filepath.Join(appDir, "modlists", modlistId, image)
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	base64Image := "data:image/png;base64," + base64.StdEncoding.EncodeToString(data)
	return base64Image, nil
}

func DeleteModlist(ctx context.Context, db *sql.DB, modlistId string) error {
	_, err := db.ExecContext(ctx, `DELETE FROM modlists WHERE id = ?`, modlistId)
	if err != nil {
		return fmt.Errorf("failed to delete modlist: %w", err)
	}

	go func() {
		appDir, err := utils.GetAppDir()
		if err != nil {
			log.Printf("Failed to get app directory for cleanup: %v", err)
			return
		}
		path := filepath.Join(appDir, "modlists", modlistId)
		if err := os.RemoveAll(path); err != nil {
			log.Printf("Failed to delete modlist directory %s: %v", path, err)
		}
	}()

	return nil
}
