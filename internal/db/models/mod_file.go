package models

type ModFile struct {
	ID             string `json:"id"`
	ModID          string `json:"mod_id"`
	Hash           string `json:"hash"`
	Type           string `json:"type"`
	Path           string `json:"path"`
	SourceFilePath string `json:"source_file_path,omitempty"`
	PatchFilePath  string `json:"patch_file_path,omitempty"`
	BsaFiles       string `json:"bsa_files,omitempty"`
}
