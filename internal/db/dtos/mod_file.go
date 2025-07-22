package dtos

type ModFileDTO struct {
	ID             string  `json:"id"`
	Hash           string  `json:"hash"`
	Type           string  `json:"type"`
	Path           string  `json:"path"`
	SourceFilePath *string `json:"source_file_path"`
	PatchFilePath  *string `json:"patch_file_path"`
	BsaFiles       *string `json:"bsa_files"`
	Size           int64   `json:"size"`
}
