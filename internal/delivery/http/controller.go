package httpdelivery

import (
	"errors"
	"net/http"

	"github.com/Nonameipal/AnalogYouTube/internal/errs"
	"github.com/Nonameipal/AnalogYouTube/internal/usecase/ports"
)

type Controller struct {
	service ports.ServiceI
}

func NewController(svc ports.ServiceI) *Controller {
	return &Controller{service: svc}
}

func (ctrl *Controller) handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, errs.ErrNotFound) ||
		errors.Is(err, errs.ErrUserNotFound) ||
		errors.Is(err, errs.ErrVideoNotFound) ||
		errors.Is(err, errs.ErrCommentNotFound) ||
		errors.Is(err, errs.ErrCategoryNotFound) ||
		errors.Is(err, errs.ErrChatNotFound) ||
		errors.Is(err, errs.ErrPlaylistNotFound):
		writeJSON(w, http.StatusNotFound, CommonError{Error: err.Error()})
	case errors.Is(err, errs.ErrInvalidRequestBody) ||
		errors.Is(err, errs.ErrInvalidFieldValue) ||
		errors.Is(err, errs.ErrCannotSubscribeToYourself) ||
		errors.Is(err, errs.ErrCannotDonateToYourself) ||
		errors.Is(err, errs.ErrCannotCreateChatWithYourself):
		writeJSON(w, http.StatusBadRequest, CommonError{Error: err.Error()})
	case errors.Is(err, errs.ErrIncorrectUsernameOrPassword) || errors.Is(err, errs.ErrInvalidToken):
		writeJSON(w, http.StatusUnauthorized, CommonError{Error: err.Error()})
	case errors.Is(err, errs.ErrAccessDenied):
		writeJSON(w, http.StatusForbidden, CommonError{Error: err.Error()})
	case errors.Is(err, errs.ErrUsernameAlreadyExists) || errors.Is(err, errs.ErrEmailAlreadyExists):
		writeJSON(w, http.StatusUnprocessableEntity, CommonError{Error: err.Error()})
	default:
		writeJSON(w, http.StatusInternalServerError, CommonError{Error: err.Error()})
	}
}
