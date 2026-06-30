package controller

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Nonameipal/AnalogYouTube/internal/errs"
	"github.com/Nonameipal/AnalogYouTube/internal/models/domain"
	"github.com/Nonameipal/AnalogYouTube/internal/models/dto"
	"github.com/Nonameipal/AnalogYouTube/pkg"
)

func (ctrl *Controller) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	userID, err := getIDFromRequest(r, "id")
	if err != nil {
		ctrl.handleError(w, err)
		return
	}
	var viewerID *int
	tokenString, err := ctrl.extractTokenFromHeader(r, authorizationHeader)
	if err == nil {
		id, isRefresh, _, err := pkg.ParseToken(tokenString)
		if err != nil {
			ctrl.handleError(w, errs.ErrInvalidToken)
			return
		}

		if isRefresh {
			ctrl.handleError(w, errs.ErrInvalidToken)
			return
		}

		viewerID = &id
	}	
	profile, err := ctrl.service.GetUserProfile(userID, viewerID)
	if err != nil {
		ctrl.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, profile)
}

func (ctrl *Controller) GetUserVideos(w http.ResponseWriter, r *http.Request) {
	userID, err := getIDFromRequest(r, "id")
	if err != nil {
		ctrl.handleError(w, err)
		return
	}

	videos, err := ctrl.service.GetUserVideos(userID)
	if err != nil {
		ctrl.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, videos)
}

func (ctrl *Controller) UpdateMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		ctrl.handleError(w, errs.ErrInvalidToken)
		return
	}

	var input dto.UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		ctrl.handleError(w, errors.Join(errs.ErrInvalidRequestBody, err))
		return
	}

	user, err := ctrl.service.UpdateUserProfile(userID, domain.User{
		Username:    input.Username,
		Email:       input.Email,
		AvatarURL:   input.AvatarURL,
		Description: input.Description,
	})
	if err != nil {
		ctrl.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, user)
}
