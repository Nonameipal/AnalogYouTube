package httpdelivery

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Nonameipal/AnalogYouTube/internal/delivery/http/dto"
	"github.com/Nonameipal/AnalogYouTube/internal/domain"
	"github.com/Nonameipal/AnalogYouTube/internal/errs"
)

// CreateComment godoc
// @Summary Создать комментарий
// @Description Добавляет комментарий к видео. Для ответа на комментарий передайте parent_id.
// @Tags Comments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID видео"
// @Param input body dto.CreateCommentRequest true "Текст комментария. parent_id  необязательный, указывает ID родительского комментария"
// @Success 201 {object} domain.Comment
// @Failure 400 {object} CommonError
// @Failure 401 {object} CommonError
// @Router /api/videos/{id}/comments [post]
func (h *Handler) CreateComment(w http.ResponseWriter, r *http.Request) {
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

	var input dto.CreateCommentRequest
	if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.handleError(w, errors.Join(errs.ErrInvalidRequestBody, err))
		return
	}

	comment, err := h.service.CreateComment(userID, videoID, domain.Comment{
		Text:     input.CommentText(),
		ParentID: input.ParentID,
	})
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, comment)
}

// GetVideoComments godoc
// @Summary Комментарии видео
// @Description Список комментариев под видео.
// @Tags Comments
// @Produce json
// @Param id path int true "ID видео"
// @Success 200 {array} domain.Comment
// @Failure 400 {object} CommonError
// @Router /api/videos/{id}/comments [get]
func (h *Handler) GetVideoComments(w http.ResponseWriter, r *http.Request) {
	videoID, err := getIDFromRequest(r, "id")
	if err != nil {
		h.handleError(w, err)
		return
	}

	comments, err := h.service.GetVideoComments(videoID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, comments)
}

// UpdateComment godoc
// @Summary Обновить комментарий
// @Description Меняет текст комментария.
// @Tags Comments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID комментария"
// @Param input body dto.UpdateCommentRequest true "Новый текст"
// @Success 200 {object} domain.Comment
// @Failure 400 {object} CommonError
// @Failure 401 {object} CommonError
// @Failure 403 {object} CommonError
// @Router /api/comments/{id} [put]
func (h *Handler) UpdateComment(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		h.handleError(w, errs.ErrInvalidToken)
		return
	}

	userRole, ok := getUserRoleFromContext(r)
	if !ok {
		h.handleError(w, errs.ErrInvalidToken)
		return
	}

	commentID, err := getIDFromRequest(r, "id")
	if err != nil {
		h.handleError(w, err)
		return
	}

	var input dto.UpdateCommentRequest
	if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.handleError(w, errors.Join(errs.ErrInvalidRequestBody, err))
		return
	}

	comment, err := h.service.UpdateComment(userID, userRole, domain.Comment{ID: commentID, Text: input.CommentText()})
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, comment)
}

// DeleteComment godoc
// @Summary Удалить комментарий
// @Description Удаляет комментарий.
// @Tags Comments
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID комментария"
// @Success 200 {object} CommonResponse
// @Failure 401 {object} CommonError
// @Failure 403 {object} CommonError
// @Router /api/comments/{id} [delete]
func (h *Handler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		h.handleError(w, errs.ErrInvalidToken)
		return
	}

	userRole, ok := getUserRoleFromContext(r)
	if !ok {
		h.handleError(w, errs.ErrInvalidToken)
		return
	}

	commentID, err := getIDFromRequest(r, "id")
	if err != nil {
		h.handleError(w, err)
		return
	}

	if err = h.service.DeleteComment(userID, userRole, commentID); err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, CommonResponse{Message: "Comment deleted successfully"})
}
