package dto

type CreateVideoRequest struct {
	Title string `json:"title"`
	Description string `json:"description"`
	VideoURL string `json:"video_url"`
	ThumbnailURL string `json:"thumbnail_url"`
	CategoryName *string `json:"category_name"`
}

type UpdateVideoRequest struct {
	Title string `json:"title"`
	Description string `json:"description"`
	VideoURL string `json:"video_url"`
	ThumbnailURL string `json:"thumbnail_url"`
	CategoryName *string `json:"category_name"`
}
