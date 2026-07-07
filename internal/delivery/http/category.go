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

// CreateCategory godoc
// @Summary Создать категорию
// @Description Новая категория для видео.
// @Tags Categories,Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body dto.CreateCategoryRequest true "Данные категории"
// @Success 201 {object} domain.Category
// @Failure 400 {object} CommonError
// @Failure 401 {object} CommonError
// @Failure 403 {object} CommonError
// @Router /api/categories [post]
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

// GetAllCategories godoc
// @Summary Все категории
// @Description Список категорий.
// @Tags Categories
// @Produce json
// @Success 200 {array} domain.Category
// @Router /api/categories [get]
func (h *Handler) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.service.GetAllCategories()
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, categories)
}

// GetCategoryByName godoc
// @Summary Категория по имени
// @Description Данные одной категории.
// @Tags Categories
// @Produce json
// @Param name path string true "Имя категории"
// @Success 200 {object} domain.Category
// @Failure 404 {object} CommonError
// @Router /api/categories/{name} [get]
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

// UpdateCategory godoc
// @Summary Обновить категорию
// @Description Меняет название или описание категории.
// @Tags Categories,Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID категории"
// @Param input body dto.UpdateCategoryRequest true "Новые данные"
// @Success 200 {object} domain.Category
// @Failure 400 {object} CommonError
// @Failure 401 {object} CommonError
// @Failure 403 {object} CommonError
// @Router /api/categories/{id} [put]
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

// DeleteCategory godoc
// @Summary Удалить категорию
// @Description Удаляет категорию.
// @Tags Categories,Admin
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID категории"
// @Success 200 {object} CommonResponse
// @Failure 401 {object} CommonError
// @Failure 403 {object} CommonError
// @Router /api/categories/{id} [delete]
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
