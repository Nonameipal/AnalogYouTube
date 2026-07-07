package dto

type UpdateProfileRequest struct {
	Username    string `json:"username"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Description string `json:"description"`
}
