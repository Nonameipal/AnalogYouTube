package repository

import (
	dbModels "github.com/Nonameipal/AnalogYouTube/internal/models/db"
	"github.com/Nonameipal/AnalogYouTube/internal/models/domain"
)

func (r *Repository) CreateUser(user domain.User) error {
	_, err := r.db.Exec(`
		INSERT INTO users (full_name, username, email, password, role, avatar_url)
		VALUES ($1, $2, NULLIF($3, ''), $4, $5, NULLIF($6, ''))`,
		user.FullName,
		user.Username,
		user.Email,
		user.Password,
		user.Role,
		user.AvatarURL,
	)
	if err != nil {
		return r.translateError(err)
	}

	return nil
}

func (r *Repository) GetUserByUsername(username string) (domain.User, error) {
	var dbUser dbModels.User
	if err := r.db.Get(&dbUser, `
		SELECT id, full_name, username, email, password, role, avatar_url, created_at, updated_at
		FROM users
		WHERE username = $1`, username); err != nil {
		return domain.User{}, r.translateError(err)
	}

	return dbUser.ToDomain(), nil
}

func (r *Repository) GetUserByEmail(email string) (domain.User, error) {
	var dbUser dbModels.User
	if err := r.db.Get(&dbUser, `
		SELECT id, full_name, username, email, password, role, avatar_url, created_at, updated_at
		FROM users
		WHERE email = $1`, email); err != nil {
		return domain.User{}, r.translateError(err)
	}

	return dbUser.ToDomain(), nil
}

func (r *Repository) GetUserByID(id int) (domain.User, error) {
	var dbUser dbModels.User
	if err := r.db.Get(&dbUser, `
		SELECT id, full_name, username, email, password, role, avatar_url, created_at, updated_at
		FROM users
		WHERE id = $1`, id); err != nil {
		return domain.User{}, r.translateError(err)
	}

	return dbUser.ToDomain(), nil
}
