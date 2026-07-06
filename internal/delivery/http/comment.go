package httpdelivery

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Nonameipal/AnalogYouTube/internal/delivery/http/dto"
	"github.com/Nonameipal/AnalogYouTube/internal/domain"
	"github.com/Nonameipal/AnalogYouTube/internal/errs"
)

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

	comment, err := h.service.CreateComment(userID, videoID, domain.Comment{Text: input.CommentText()})
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, comment)
}

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
