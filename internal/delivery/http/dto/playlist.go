package dto

type CreatePlaylistRequest struct {
	Name string `json:"name"`
	Description string `json:"description"`
}

type UpdatePlaylistRequest struct {
	Name string `json:"name"`
	Description string `json:"description"`
}

type AddVideoToPlaylistRequest struct {
	VideoID int `json:"video_id"`
}
