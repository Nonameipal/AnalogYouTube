package httpdelivery

import (
	"net/http"

	"github.com/Nonameipal/AnalogYouTube/internal/errs"
)

// LikeVideo godoc
// @Summary Поставить лайк
// @Description Лайк от текущего пользователя.
// @Tags Likes
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID видео"
// @Success 200 {object} CommonResponse
// @Failure 401 {object} CommonError
// @Router /api/videos/{id}/like [post]
func (h *Handler) LikeVideo(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		h.handleError(w, errs.ErrInvalidToken)
		return
	}

	videoID, err := getIDFromRequest(r, "id")
	if err != nil {
		h.handleError(w, err)
		return
	}

	if err = h.service.LikeVideo(userID, videoID); err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, CommonResponse{Message: "Video liked successfully"})
}

// UnlikeVideo godoc
// @Summary Убрать лайк
// @Description Удаляет мой лайк с видео.
// @Tags Likes
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID видео"
// @Success 200 {object} CommonResponse
// @Failure 401 {object} CommonError
// @Router /api/videos/{id}/like [delete]
func (h *Handler) UnlikeVideo(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		h.handleError(w, errs.ErrInvalidToken)
		return
	}

	videoID, err := getIDFromRequest(r, "id")
	if err != nil {
		h.handleError(w, err)
		return
	}

	if err = h.service.UnlikeVideo(userID, videoID); err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, CommonResponse{Message: "Video unliked successfully"})
}

// GetVideoLikesCount godoc
// @Summary Количество лайков
// @Description Сколько лайков у видео.
// @Tags Likes
// @Produce json
// @Param id path int true "ID видео"
// @Success 200 {object} map[string]int
// @Failure 400 {object} CommonError
// @Router /api/videos/{id}/likes/count [get]
func (h *Handler) GetVideoLikesCount(w http.ResponseWriter, r *http.Request) {
	videoID, err := getIDFromRequest(r, "id")
	if err != nil {
		h.handleError(w, err)
		return
	}

	count, err := h.service.GetVideoLikesCount(videoID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]int{"likes_count": count})
}

// IsVideoLikedByMe godoc
// @Summary Мой лайк
// @Description Показывает, лайкнул ли я видео.
// @Tags Likes
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID видео"
// @Success 200 {object} map[string]bool
// @Failure 401 {object} CommonError
// @Router /api/videos/{id}/liked [get]
func (h *Handler) IsVideoLikedByMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		h.handleError(w, errs.ErrInvalidToken)
		return
	}

	videoID, err := getIDFromRequest(r, "id")
	if err != nil {
		h.handleError(w, err)
		return
	}

	liked, err := h.service.IsVideoLikedByUser(userID, videoID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]bool{"liked": liked})
}
