package httpdelivery

import (
	"net/http"

	"github.com/Nonameipal/AnalogYouTube/internal/domain"
	"github.com/Nonameipal/AnalogYouTube/internal/errs"
	"github.com/Nonameipal/AnalogYouTube/pkg"
)

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
