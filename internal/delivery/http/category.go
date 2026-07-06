package httpdelivery

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Nonameipal/AnalogYouTube/internal/delivery/http/dto"
	"github.com/Nonameipal/AnalogYouTube/internal/domain"
	"github.com/Nonameipal/AnalogYouTube/internal/errs"
	"github.com/gorilla/mux"
)

func (h *Handler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var input dto.CreateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.handleError(w, errors.Join(errs.ErrInvalidRequestBody, err))
		return
	}

	category, err := h.service.CreateCategory(domain.Category{
		Name:        input.Name,
		Description: input.Description,
	})
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, category)
}

func (h *Handler) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.service.GetAllCategories()
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, categories)
}

func (h *Handler) GetCategoryByName(w http.ResponseWriter, r *http.Request) {
	categoryName := mux.Vars(r)["name"]
	if categoryName == "" {
		h.handleError(w, errs.ErrInvalidFieldValue)
		return
	}

	category, err := h.service.GetCategoryByName(categoryName)
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, category)
}

func (h *Handler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	categoryID, err := getIDFromRequest(r, "id")
	if err != nil {
		h.handleError(w, err)
		return
	}

	var input dto.UpdateCategoryRequest
	if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.handleError(w, errors.Join(errs.ErrInvalidRequestBody, err))
		return
	}

	category, err := h.service.UpdateCategory(domain.Category{
		ID:          categoryID,
		Name:        input.Name,
		Description: input.Description,
	})
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, category)
}

func (h *Handler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	categoryID, err := getIDFromRequest(r, "id")
	if err != nil {
		h.handleError(w, err)
		return
	}

	if err = h.service.DeleteCategory(categoryID); err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, CommonResponse{Message: "Category deleted successfully"})
}
