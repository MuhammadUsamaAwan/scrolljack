package models

type ModArchive struct {
	ID            string `db:"id"`
	ModID         string `db:"mod_id"`
	Type          string `db:"type"`
	NexusGameName string `db:"nexus_game_name"`
	NexusModID    int    `db:"nexus_mod_id"`
	NexusFileID   int    `db:"nexus_file_id"`
	DirectURL     string `db:"direct_url"`
	Version       string `db:"version"`
	Size          int    `db:"size"`
	Description   string `db:"description"`
}
