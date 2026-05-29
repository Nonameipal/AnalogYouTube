package controller

import (
	"errors"
	"net/http"

	"github.com/Nonameipal/AnalogYouTube/internal/contracts"
	"github.com/Nonameipal/AnalogYouTube/internal/errs"
)

type Controller struct {
	service contracts.ServiceI
}

func NewController(svc contracts.ServiceI) *Controller {
	return &Controller{service: svc}
}

func (ctrl *Controller) handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, errs.ErrNotFound) || errors.Is(err, errs.ErrUserNotFound):
		writeJSON(w, http.StatusNotFound, CommonError{Error: err.Error()})
	case errors.Is(err, errs.ErrInvalidRequestBody) || errors.Is(err, errs.ErrInvalidFieldValue):
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
