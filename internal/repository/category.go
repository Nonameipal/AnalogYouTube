package repository

import (
	"github.com/Nonameipal/AnalogYouTube/internal/errs"
	dbModels "github.com/Nonameipal/AnalogYouTube/internal/models/db"
	"github.com/Nonameipal/AnalogYouTube/internal/models/domain"
)

func (r *Repository) CreateCategory(category domain.Category) (domain.Category, error) {
	var dbCategory dbModels.Category
	err := r.db.Get(&dbCategory, 
		`INSERT INTO categories (name, description)
		VALUES ($1, NULLIF($2, ''))
		RETURNING id, name, description, created_at, updated_at`,
		category.Name,
		category.Description,
	)
	if err != nil {
		return domain.Category{}, r.translateError(err)
	}

	return dbCategory.ToDomain(), nil
}

func (r *Repository) GetAllCategories() ([]domain.Category, error) {
	var dbCategories []dbModels.Category
	err := r.db.Select(&dbCategories, 
		`SELECT id, name, description, created_at, updated_at
		FROM categories
		ORDER BY name ASC`)
	if err != nil {
		return nil, r.translateError(err)
	}

	categories := make([]domain.Category, 0, len(dbCategories))
	for _, category := range dbCategories {
		categories = append(categories, category.ToDomain())
	}

	return categories, nil
}

func (r *Repository) GetCategoryByID(id int) (domain.Category, error) {
	var dbCategory dbModels.Category
	if err := r.db.Get(&dbCategory,
		`SELECT id, name, description, created_at, updated_at
		FROM categories
		WHERE id = $1`, id); err != nil {
		return domain.Category{}, r.translateError(err)
	}

	return dbCategory.ToDomain(), nil
}

func (r *Repository) UpdateCategory(category domain.Category) (domain.Category, error) {
	var dbCategory dbModels.Category
	err := r.db.Get(&dbCategory, 
		`UPDATE categories
		SET name = $1,
		    description = NULLIF($2, ''),
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $3
		RETURNING id, name, description, created_at, updated_at`,
		category.Name,
		category.Description,
		category.ID,
	)
	if err != nil {
		return domain.Category{}, r.translateError(err)
	}

	return dbCategory.ToDomain(), nil
}

func (r *Repository) DeleteCategory(id int) error {
	result, err := r.db.Exec(
		`DELETE FROM categories
		WHERE id = $1`, id)
	if err != nil {
		return r.translateError(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return r.translateError(err)
	}
	if rowsAffected == 0 {
		return errs.ErrNotFound
	}

	return nil
}
