package models

import "database/sql"

type ModArchive struct {
	ID            string         `db:"id"`
	ModID         string         `db:"mod_id"`
	Type          string         `db:"type"`
	NexusGameName sql.NullString `db:"nexus_game_name"`
	NexusModID    sql.NullInt64  `db:"nexus_mod_id"`
	NexusFileID   sql.NullInt64  `db:"nexus_file_id"`
	DirectURL     sql.NullString `db:"direct_url"`
	Version       sql.NullString `db:"version"`
	Size          sql.NullInt64  `db:"size"`
	Description   sql.NullString `db:"description"`
}
