package db

import (
	"database/sql"
	"time"

	"github.com/Nonameipal/AnalogYouTube/internal/models/domain"
)

type User struct {
	ID int `db:"id"`
	Username string `db:"username"`
	Email sql.NullString `db:"email"`
	Password string `db:"password"`
	Role string `db:"role"`
	AvatarURL sql.NullString `db:"avatar_url"`
	Description sql.NullString `db:"description"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (u User) ToDomain() domain.User {
	return domain.User{
		ID: u.ID,
		Username: u.Username,
		Email: u.Email.String,
		Password: u.Password,
		Role: u.Role,
		AvatarURL: u.AvatarURL.String,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		Description: u.Description.String,
	}
}
