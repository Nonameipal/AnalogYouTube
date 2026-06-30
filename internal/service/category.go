package service

import (
	"errors"
	"strings"

	"github.com/Nonameipal/AnalogYouTube/internal/errs"
	"github.com/Nonameipal/AnalogYouTube/internal/models/domain"
)

func (s *Service) CreateCategory(category domain.Category) (domain.Category, error) {
	category.Name = strings.TrimSpace(category.Name)
	category.Description = strings.TrimSpace(category.Description)

	if category.Name == "" {
		return domain.Category{}, errs.ErrInvalidFieldValue
	}

	return s.repository.CreateCategory(category)
}

func (s *Service) GetAllCategories() ([]domain.Category, error) {
	return s.repository.GetAllCategories()
}

func (s *Service) GetCategoryByName(name string) (domain.Category, error) {
	if name == "" {
		return domain.Category{}, errs.ErrInvalidFieldValue
	}

	category, err := s.repository.GetCategoryByName(name)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return domain.Category{}, errs.ErrCategoryNotFound
		}
		return domain.Category{}, err
	}

	return category, nil
}

func (s *Service) UpdateCategory(category domain.Category) (domain.Category, error) {
	category.Name = strings.TrimSpace(category.Name)
	category.Description = strings.TrimSpace(category.Description)

	if category.ID <= 0 || category.Name == "" {
		return domain.Category{}, errs.ErrInvalidFieldValue
	}

	category, err := s.repository.UpdateCategory(category)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return domain.Category{}, errs.ErrCategoryNotFound
		}
		return domain.Category{}, err
	}

	return category, nil
}

func (s *Service) DeleteCategory(id int) error {
	if id <= 0 {
		return errs.ErrInvalidFieldValue
	}

	if err := s.repository.DeleteCategory(id); err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return errs.ErrCategoryNotFound
		}
		return err
	}

	return nil
}
