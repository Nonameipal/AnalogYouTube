package httpdelivery

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Nonameipal/AnalogYouTube/internal/delivery/http/dto"
	"github.com/Nonameipal/AnalogYouTube/internal/domain"
	"github.com/Nonameipal/AnalogYouTube/internal/errs"
)

func (ctrl *Controller) CreateComment(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		ctrl.handleError(w, errs.ErrInvalidToken)
		return
	}

	videoID, err := getIDFromRequest(r, "id")
	if err != nil {
		ctrl.handleError(w, err)
		return
	}

	var input dto.CreateCommentRequest
	if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
		ctrl.handleError(w, errors.Join(errs.ErrInvalidRequestBody, err))
		return
	}

	comment, err := ctrl.service.CreateComment(userID, videoID, domain.Comment{Text: input.Text})
	if err != nil {
		ctrl.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, comment)
}

func (ctrl *Controller) GetVideoComments(w http.ResponseWriter, r *http.Request) {
	videoID, err := getIDFromRequest(r, "id")
	if err != nil {
		ctrl.handleError(w, err)
		return
	}

	comments, err := ctrl.service.GetVideoComments(videoID)
	if err != nil {
		ctrl.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, comments)
}

func (ctrl *Controller) UpdateComment(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		ctrl.handleError(w, errs.ErrInvalidToken)
		return
	}

	userRole, ok := getUserRoleFromContext(r)
	if !ok {
		ctrl.handleError(w, errs.ErrInvalidToken)
		return
	}

	commentID, err := getIDFromRequest(r, "id")
	if err != nil {
		ctrl.handleError(w, err)
		return
	}

	var input dto.UpdateCommentRequest
	if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
		ctrl.handleError(w, errors.Join(errs.ErrInvalidRequestBody, err))
		return
	}

	comment, err := ctrl.service.UpdateComment(userID, userRole, domain.Comment{ID: commentID, Text: input.Text})
	if err != nil {
		ctrl.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, comment)
}

func (ctrl *Controller) DeleteComment(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		ctrl.handleError(w, errs.ErrInvalidToken)
		return
	}

	userRole, ok := getUserRoleFromContext(r)
	if !ok {
		ctrl.handleError(w, errs.ErrInvalidToken)
		return
	}

	commentID, err := getIDFromRequest(r, "id")
	if err != nil {
		ctrl.handleError(w, err)
		return
	}

	if err = ctrl.service.DeleteComment(userID, userRole, commentID); err != nil {
		ctrl.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, CommonResponse{Message: "Comment deleted successfully"})
}
