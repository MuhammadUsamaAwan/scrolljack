package dtos

type ModArchiveDTO struct {
	ID            string  `json:"id"`
	Hash          string  `json:"hash"`
	Type          string  `json:"type"`
	NexusGameName *string `json:"nexus_game_name"`
	NexusModID    *string `json:"nexus_mod_id"`
	NexusFileID   *string `json:"nexus_file_id"`
	DirectURL     *string `json:"direct_url"`
	Version       *string `json:"version"`
	Size          *int64  `json:"size"`
	Description   *string `json:"description"`
}
