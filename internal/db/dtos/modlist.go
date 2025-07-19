package dtos

type ModlistDTO struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Author      string `json:"author"`
	Description string `json:"description"`
	Image       string `json:"image"`
	GameType    string `json:"game_type"`
	Version     string `json:"version"`
	IsNSFW      bool   `json:"is_nsfw"`
	CreatedAt   string `json:"created_at"`
}
