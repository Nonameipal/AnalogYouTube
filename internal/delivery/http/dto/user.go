package dto

type UpdateProfileRequest struct {
	Username string `json:"username"`
	Email string `json:"email"`
	AvatarURL string `json:"avatar_url"`
	Description string `json:"description"`
}
