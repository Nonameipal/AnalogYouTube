package repository

import (
	"context"

	"github.com/Nonameipal/AnalogYouTube/internal/errs"
	dbModels "github.com/Nonameipal/AnalogYouTube/internal/models/db"
	"github.com/Nonameipal/AnalogYouTube/internal/models/domain"
)

func (r *Repository) CreateCategory(category domain.Category) (domain.Category, error) {
	ctx := context.Background()
	var dbCategory dbModels.Category
	err := r.db.QueryRow(ctx,
		`INSERT INTO categories (name, description)
		VALUES ($1, NULLIF($2, ''))
		RETURNING id, name, description, created_at, updated_at`,
		category.Name,
		category.Description,
	).Scan(&dbCategory.ID, &dbCategory.Name, &dbCategory.Description, &dbCategory.CreatedAt, &dbCategory.UpdatedAt)
	if err != nil {
		return domain.Category{}, r.translateError(err)
	}

	return dbCategory.ToDomain(), nil
}

func (r *Repository) GetAllCategories() ([]domain.Category, error) {
	ctx := context.Background()
	rows, err := r.db.Query(ctx,
		`SELECT id, name, description, created_at, updated_at
		FROM categories
		ORDER BY name ASC`)
	if err != nil {
		return nil, r.translateError(err)
	}
	defer rows.Close()

	var dbCategories []dbModels.Category
	for rows.Next() {
		var category dbModels.Category
		if err := rows.Scan(&category.ID, &category.Name, &category.Description, &category.CreatedAt, &category.UpdatedAt); err != nil {
			return nil, r.translateError(err)
		}
		dbCategories = append(dbCategories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, r.translateError(err)
	}

	categories := make([]domain.Category, 0, len(dbCategories))
	for _, category := range dbCategories {
		categories = append(categories, category.ToDomain())
	}

	return categories, nil
}

func (r *Repository) GetCategoryByName(name string) (domain.Category, error) {
	ctx := context.Background()
	var dbCategory dbModels.Category
	err := r.db.QueryRow(ctx,
		`SELECT id, name, description, created_at, updated_at
		FROM categories
		WHERE LOWER(name) = LOWER($1) 
		LIMIT 1`, name).Scan(&dbCategory.ID, &dbCategory.Name, &dbCategory.Description, &dbCategory.CreatedAt, &dbCategory.UpdatedAt)
	if err != nil {
		return domain.Category{}, r.translateError(err)
	}

	return dbCategory.ToDomain(), nil
}

func (r *Repository) UpdateCategory(category domain.Category) (domain.Category, error) {
	ctx := context.Background()
	var dbCategory dbModels.Category
	err := r.db.QueryRow(ctx,
		`UPDATE categories
		SET name = $1,
		    description = NULLIF($2, ''),
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $3
		RETURNING id, name, description, created_at, updated_at`,
		category.Name,
		category.Description,
		category.ID,
	).Scan(&dbCategory.ID, &dbCategory.Name, &dbCategory.Description, &dbCategory.CreatedAt, &dbCategory.UpdatedAt)
	if err != nil {
		return domain.Category{}, r.translateError(err)
	}

	return dbCategory.ToDomain(), nil
}

func (r *Repository) DeleteCategory(id int) error {
	ctx := context.Background()
	result, err := r.db.Exec(ctx,
		`DELETE FROM categories
		WHERE id = $1`, id)
	if err != nil {
		return r.translateError(err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return errs.ErrNotFound
	}

	return nil
}
