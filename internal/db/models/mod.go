package models

type Mod struct {
	ID          string `json:"id"`
	ProfileID   string `json:"profile_id"`
	Name        string `json:"name"`
	IsSeparator bool   `json:"is_separator"`
	Order       int    `json:"order"`
	ModOrder    int    `json:"mod_order"`
	IsActive    bool   `json:"is_active"`
}
