package usecase

import (
	"errors"
	"strings"

	"github.com/Nonameipal/AnalogYouTube/internal/domain"
	"github.com/Nonameipal/AnalogYouTube/internal/errs"
)

func (uc *Usecase) CreateCategory(category domain.Category) (domain.Category, error) {
	category.Name = strings.TrimSpace(category.Name)
	category.Description = strings.TrimSpace(category.Description)

	if category.Name == "" {
		return domain.Category{}, errs.ErrInvalidFieldValue
	}

	return uc.repository.CreateCategory(category)
}

func (uc *Usecase) GetAllCategories() ([]domain.Category, error) {
	return uc.repository.GetAllCategories()
}

func (uc *Usecase) GetCategoryByName(name string) (domain.Category, error) {
	if name == "" {
		return domain.Category{}, errs.ErrInvalidFieldValue
	}

	category, err := uc.repository.GetCategoryByName(name)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return domain.Category{}, errs.ErrCategoryNotFound
		}
		return domain.Category{}, err
	}

	return category, nil
}

func (uc *Usecase) UpdateCategory(category domain.Category) (domain.Category, error) {
	category.Name = strings.TrimSpace(category.Name)
	category.Description = strings.TrimSpace(category.Description)

	if category.ID <= 0 || category.Name == "" {
		return domain.Category{}, errs.ErrInvalidFieldValue
	}

	category, err := uc.repository.UpdateCategory(category)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return domain.Category{}, errs.ErrCategoryNotFound
		}
		return domain.Category{}, err
	}

	return category, nil
}

func (uc *Usecase) DeleteCategory(id int) error {
	if id <= 0 {
		return errs.ErrInvalidFieldValue
	}

	if err := uc.repository.DeleteCategory(id); err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return errs.ErrCategoryNotFound
		}
		return err
	}

	return nil
}
