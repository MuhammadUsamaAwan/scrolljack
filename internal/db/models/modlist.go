package models

import "database/sql"

type Modlist struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Author      sql.NullString `json:"author"`
	Description sql.NullString `json:"description"`
	Website     sql.NullString `json:"website"`
	Image       sql.NullString `json:"image"`
	Readme      sql.NullString `json:"readme"`
	GameType    sql.NullString `json:"game_type"`
	Version     sql.NullString `json:"version"`
	IsNSFW      bool           `json:"is_nsfw"`
	CreatedAt   string         `json:"created_at"`
}
