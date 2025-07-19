package models

type ProfileFile struct {
	ID        string `json:"id"`
	ProfileID string `json:"profile_id"`
	Name      string `json:"name"`
	FilePath  string `json:"file_path"`
}
