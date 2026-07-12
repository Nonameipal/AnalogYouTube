package usecase

import (
	"errors"
	"testing"

	"github.com/Nonameipal/AnalogYouTube/internal/domain"
	"github.com/Nonameipal/AnalogYouTube/internal/errs"
)

func TestCategoryUsecase(t *testing.T) {
	t.Run("create trims and saves", func(t *testing.T) {
		repo := newFakeRepository()
		uc := newTestUsecase(repo)

		category, err := uc.CreateCategory(domain.Category{Name: " Programming ", Description: " Code "})
		if err != nil {
			t.Fatalf("CreateCategory returned error: %v", err)
		}
		if category.Name != "Programming" || category.Description != "Code" {
			t.Fatalf("category was not trimmed: %+v", category)
		}
	})

	t.Run("not found is mapped", func(t *testing.T) {
		uc := newTestUsecase(newFakeRepository())

		_, err := uc.GetCategoryByName("missing")
		if !errors.Is(err, errs.ErrCategoryNotFound) {
			t.Fatalf("expected ErrCategoryNotFound, got %v", err)
		}
	})

	t.Run("update validates id and name", func(t *testing.T) {
		uc := newTestUsecase(newFakeRepository())

		_, err := uc.UpdateCategory(domain.Category{ID: 0, Name: "Programming"})
		if !errors.Is(err, errs.ErrInvalidFieldValue) {
			t.Fatalf("expected ErrInvalidFieldValue, got %v", err)
		}
	})

	t.Run("delete maps not found", func(t *testing.T) {
		repo := newFakeRepository()
		repo.deleteCategoryErr = errs.ErrNotFound
		uc := newTestUsecase(repo)

		err := uc.DeleteCategory(1)
		if !errors.Is(err, errs.ErrCategoryNotFound) {
			t.Fatalf("expected ErrCategoryNotFound, got %v", err)
		}
	})
}
