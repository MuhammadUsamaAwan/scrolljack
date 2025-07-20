package dtos

type Mod struct {
	ID          string `json:"id"`
	ProfileID   string `json:"profile_id"`
	Name        string `json:"name"`
	Order       int    `json:"order"`
	IsSeparator bool   `json:"is_separator"`
}

type GroupedMod struct {
	Separator string `json:"separator"`
	Mods      []Mod  `json:"mods"`
}
