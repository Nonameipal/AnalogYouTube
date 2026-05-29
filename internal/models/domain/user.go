package domain

import "time"

const (
	UserRole  = "USER"
	AdminRole = "ADMIN"
)

type User struct {
	ID        int       `json:"id"`
	FullName  string    `json:"full_name"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Role      string    `json:"role"`
	AvatarURL string    `json:"avatar_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
