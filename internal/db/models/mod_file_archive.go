package models

type ModFileArchive struct {
	ModlistId    string `db:"modlist_id"`
	ModFileId    string `db:"mod_file_id"`
	ModArchiveId string `db:"mod_archive_id"`
}
