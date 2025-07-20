package dtos

type ModDTO struct {
	ID          string `json:"id"`
	ProfileID   string `json:"profile_id"`
	Name        string `json:"name"`
	Order       int    `json:"order"`
	ModOrder    int    `json:"mod_order"`
	IsActive    bool   `json:"is_active"`
	IsSeparator bool   `json:"is_separator"`
}

type GroupedModDTO struct {
	Separator string   `json:"separator"`
	Mods      []ModDTO `json:"mods"`
}
