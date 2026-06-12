package controller

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Nonameipal/AnalogYouTube/internal/errs"
	"github.com/Nonameipal/AnalogYouTube/internal/models/domain"
	"github.com/Nonameipal/AnalogYouTube/internal/models/dto"
)

func (ctrl *Controller) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var input dto.CreateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		ctrl.handleError(w, errors.Join(errs.ErrInvalidRequestBody, err))
		return
	}

	category, err := ctrl.service.CreateCategory(domain.Category{
		Name:        input.Name,
		Description: input.Description,
	})
	if err != nil {
		ctrl.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, category)
}

func (ctrl *Controller) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := ctrl.service.GetAllCategories()
	if err != nil {
		ctrl.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, categories)
}

func (ctrl *Controller) GetCategoryByID(w http.ResponseWriter, r *http.Request) {
	categoryID, err := getIDFromRequest(r, "id")
	if err != nil {
		ctrl.handleError(w, err)
		return
	}

	category, err := ctrl.service.GetCategoryByID(categoryID)
	if err != nil {
		ctrl.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, category)
}

func (ctrl *Controller) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	categoryID, err := getIDFromRequest(r, "id")
	if err != nil {
		ctrl.handleError(w, err)
		return
	}

	var input dto.UpdateCategoryRequest
	if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
		ctrl.handleError(w, errors.Join(errs.ErrInvalidRequestBody, err))
		return
	}

	category, err := ctrl.service.UpdateCategory(domain.Category{
		ID: categoryID,
		Name: input.Name,
		Description: input.Description,
	})
	if err != nil {
		ctrl.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, category)
}

func (ctrl *Controller) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	categoryID, err := getIDFromRequest(r, "id")
	if err != nil {
		ctrl.handleError(w, err)
		return
	}

	if err = ctrl.service.DeleteCategory(categoryID); err != nil {
		ctrl.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, CommonResponse{Message: "Category deleted successfully"})
}
