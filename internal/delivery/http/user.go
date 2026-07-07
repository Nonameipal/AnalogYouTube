package httpdelivery

import (
	"net/http"

	"github.com/Nonameipal/AnalogYouTube/internal/domain"
	"github.com/Nonameipal/AnalogYouTube/internal/errs"
	"github.com/Nonameipal/AnalogYouTube/pkg"
)

// GetUserProfile godoc
// @Summary Профиль пользователя
// @Description Публичная страница канала.
// @Tags Profile
// @Produce json
// @Param id path int true "ID пользователя"
// @Success 200 {object} domain.UserProfile
// @Failure 400 {object} CommonError
// @Failure 404 {object} CommonError
// @Router /api/users/{id} [get]
func (h *Handler) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	userID, err := getIDFromRequest(r, "id")
	if err != nil {
		h.handleError(w, err)
		return
	}
	var viewerID *int
	tokenString, err := h.extractTokenFromHeader(r, authorizationHeader)
	if err == nil {
		id, isRefresh, _, err := pkg.ParseToken(tokenString)
		if err != nil {
			h.handleError(w, errs.ErrInvalidToken)
			return
		}

		if isRefresh {
			h.handleError(w, errs.ErrInvalidToken)
			return
		}

		viewerID = &id
	}
	profile, err := h.service.GetUserProfile(userID, viewerID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, profile)
}

// GetUserVideos godoc
// @Summary Видео пользователя
// @Description Ролики выбранного автора.
// @Tags Profile
// @Produce json
// @Param id path int true "ID пользователя"
// @Success 200 {array} domain.Video
// @Failure 400 {object} CommonError
// @Failure 404 {object} CommonError
// @Router /api/users/{id}/videos [get]
func (h *Handler) GetUserVideos(w http.ResponseWriter, r *http.Request) {
	userID, err := getIDFromRequest(r, "id")
	if err != nil {
		h.handleError(w, err)
		return
	}

	videos, err := h.service.GetUserVideos(userID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, videos)
}

// UpdateMe godoc
// @Summary Обновить профиль
// @Description Меняет мои данные, аватарку или пароль.
// @Tags Profile
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param username formData string true "Новое имя пользователя"
// @Param email formData string false "Email, например user@gmail.com"
// @Param password formData string false "Новый пароль"
// @Param description formData string false "Описание канала"
// @Param avatar formData file false "Файл аватарки"
// @Success 200 {object} domain.User
// @Failure 400 {object} CommonError
// @Failure 401 {object} CommonError
// @Failure 422 {object} CommonError
// @Router /api/me [put]
func (h *Handler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		h.handleError(w, errs.ErrInvalidToken)
		return
	}

	if err := r.ParseMultipartForm(20 << 20); err != nil {
		h.handleError(w, errs.ErrInvalidRequestBody)
		return
	}

	avatarURL, err := h.saveMultipartFile(r, "avatar", "avatars", false)
	if err != nil {
		h.handleError(w, err)
		return
	}

	user, err := h.service.UpdateUserProfile(userID, domain.User{
		Username:    r.FormValue("username"),
		Email:       r.FormValue("email"),
		Password:    r.FormValue("password"),
		AvatarURL:   avatarURL,
		Description: r.FormValue("description"),
	})
	if err != nil {
		if avatarURL != "" {
			h.cleanupFailedVideoUpload(0, "", 0, "", avatarURL)
		}
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, user)
}
