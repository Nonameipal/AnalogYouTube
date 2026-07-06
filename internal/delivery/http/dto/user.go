package dto

type UpdateProfileRequest struct {
	Username    string `json:"username"`
	Email       string `json:"email"`
	Description string `json:"description"`
}
