package models

import "database/sql"

type ModFile struct {
	ID             string         `json:"id"`
	ModID          string         `json:"mod_id"`
	Hash           string         `json:"hash"`
	Type           string         `json:"type"`
	Path           string         `json:"path"`
	SourceFilePath sql.NullString `json:"source_file_path,omitempty"`
	PatchFilePath  sql.NullString `json:"patch_file_path,omitempty"`
	BsaFiles       sql.NullString `json:"bsa_files,omitempty"`
}
