package postgres

import (
	"context"

	"github.com/Nonameipal/AnalogYouTube/internal/domain"
	"github.com/Nonameipal/AnalogYouTube/internal/repository/postgres/dbmodel"
)

func (r *Repository) CreateUser(user domain.User) error {
	ctx := context.Background()
	_, err := r.db.Exec(ctx,
		`INSERT INTO users (username, email, password, role, avatar_url, description)
		VALUES ($1, $2, NULLIF($3, ''), $4, NULLIF($5, ''), NULLIF($6, ''))`,
		user.Username,
		user.Email,
		user.Password,
		user.Role,
		user.AvatarURL,
		user.Description,
	)
	if err != nil {
		return r.translateError(err)
	}

	return nil
}

func (r *Repository) GetUserByUsername(username string) (domain.User, error) {
	ctx := context.Background()
	var dbUser dbmodel.User
	err := r.db.QueryRow(ctx,
		`SELECT id, username, email, password, role, avatar_url, description, created_at, updated_at
		FROM users
		WHERE username = $1`, username).Scan(
		&dbUser.ID, &dbUser.Username, &dbUser.Email, &dbUser.Password, &dbUser.Role,
		&dbUser.AvatarURL, &dbUser.Description, &dbUser.CreatedAt, &dbUser.UpdatedAt)
	if err != nil {
		return domain.User{}, r.translateError(err)
	}

	return dbUser.ToDomain(), nil
}

func (r *Repository) GetUserByEmail(email string) (domain.User, error) {
	ctx := context.Background()
	var dbUser dbmodel.User
	err := r.db.QueryRow(ctx,
		`SELECT id, username, email, password, role, avatar_url, description, created_at, updated_at
		FROM users
		WHERE email = $1`, email).Scan(
		&dbUser.ID, &dbUser.Username, &dbUser.Email, &dbUser.Password, &dbUser.Role,
		&dbUser.AvatarURL, &dbUser.Description, &dbUser.CreatedAt, &dbUser.UpdatedAt)
	if err != nil {
		return domain.User{}, r.translateError(err)
	}

	return dbUser.ToDomain(), nil
}

func (r *Repository) GetUserByID(id int) (domain.User, error) {
	ctx := context.Background()
	var dbUser dbmodel.User
	err := r.db.QueryRow(ctx,
		`SELECT id, username, email, password, role, avatar_url, description, created_at, updated_at
		FROM users
		WHERE id = $1`, id).Scan(
		&dbUser.ID, &dbUser.Username, &dbUser.Email, &dbUser.Password, &dbUser.Role,
		&dbUser.AvatarURL, &dbUser.Description, &dbUser.CreatedAt, &dbUser.UpdatedAt)
	if err != nil {
		return domain.User{}, r.translateError(err)
	}

	return dbUser.ToDomain(), nil
}

func (r *Repository) UpdateUserProfile(user domain.User) (domain.User, error) {
	ctx := context.Background()
	var dbUser dbmodel.User
	err := r.db.QueryRow(ctx,
		`UPDATE users
		SET username = $1,
		    email = NULLIF($2, ''),
		    avatar_url = NULLIF($3, ''),
		    description = NULLIF($4, ''),
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $5
		RETURNING id, username, email, password, role, avatar_url, description, created_at, updated_at`,
		user.Username,
		user.Email,
		user.AvatarURL,
		user.Description,
		user.ID,
	).Scan(
		&dbUser.ID, &dbUser.Username, &dbUser.Email, &dbUser.Password, &dbUser.Role,
		&dbUser.AvatarURL, &dbUser.Description, &dbUser.CreatedAt, &dbUser.UpdatedAt)
	if err != nil {
		return domain.User{}, r.translateError(err)
	}

	return dbUser.ToDomain(), nil
}
