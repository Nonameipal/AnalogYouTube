package repository

import (
	dbModels "github.com/Nonameipal/AnalogYouTube/internal/models/db"
	"github.com/Nonameipal/AnalogYouTube/internal/models/domain"
)

func (r *Repository) CreateUser(user domain.User) error {
	_, err := r.db.Exec(
		`INSERT INTO users (username, email, password, role, avatar_url)
		VALUES ($1, $2, NULLIF($3, ''), $4, NULLIF($5, ''))`,
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
	if err := r.db.Get(&dbUser, 
		`SELECT id, username, email, password, role, avatar_url, created_at, updated_at
		FROM users
		WHERE username = $1`, username); err != nil {
		return domain.User{}, r.translateError(err)
	}

	return dbUser.ToDomain(), nil
}

func (r *Repository) GetUserByEmail(email string) (domain.User, error) {
	var dbUser dbModels.User
	if err := r.db.Get(&dbUser, 
		`SELECT id, username, email, password, role, avatar_url, created_at, updated_at
		FROM users
		WHERE email = $1`, email); err != nil {
		return domain.User{}, r.translateError(err)
	}

	return dbUser.ToDomain(), nil
}

func (r *Repository) GetUserByID(id int) (domain.User, error) {
	var dbUser dbModels.User
	if err := r.db.Get(&dbUser, 
		`SELECT id, username, email, password, role, avatar_url, created_at, updated_at
		FROM users
		WHERE id = $1`, id); err != nil {
		return domain.User{}, r.translateError(err)
	}

	return dbUser.ToDomain(), nil
}

func (r *Repository) UpdateUserProfile(user domain.User) (domain.User, error) {
	var dbUser dbModels.User
	err := r.db.Get(&dbUser, 
		`UPDATE users
		SET username = $1,
		    email = NULLIF($2, ''),
		    avatar_url = NULLIF($3, ''),
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $4
		RETURNING id, username, email, password, role, avatar_url, created_at, updated_at`,
		user.Username,
		user.Email,
		user.AvatarURL,
		user.ID,
	)
	if err != nil {
		return domain.User{}, r.translateError(err)
	}

	return dbUser.ToDomain(), nil
}
