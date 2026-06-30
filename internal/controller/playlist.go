package controller

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Nonameipal/AnalogYouTube/internal/errs"
	"github.com/Nonameipal/AnalogYouTube/internal/models/domain"
	"github.com/Nonameipal/AnalogYouTube/internal/models/dto"
)

func (ctrl *Controller) CreatePlaylist(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		ctrl.handleError(w, errs.ErrInvalidToken)
		return
	}

	var input dto.CreatePlaylistRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		ctrl.handleError(w, errors.Join(errs.ErrInvalidRequestBody, err))
		return
	}

	playlist, err := ctrl.service.CreatePlaylist(userID, domain.Playlist{
		Name: input.Name,
		Description: input.Description,
	})
	if err != nil {
		ctrl.handleError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, playlist)
}

func (ctrl *Controller) GetUserPlaylists(w http.ResponseWriter, r *http.Request) {
	userID, err := getIDFromRequest(r, "id")
	if err != nil {
		ctrl.handleError(w, err)
		return
	}

	playlists, err := ctrl.service.GetUserPlaylists(userID)
	if err != nil {
		ctrl.handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, playlists)
}

func (ctrl *Controller) GetPlaylistByID(w http.ResponseWriter, r *http.Request) {
	playlistID, err := getIDFromRequest(r, "id")
	if err != nil {
		ctrl.handleError(w, err)
		return
	}

	playlist, err := ctrl.service.GetPlaylistByID(playlistID)
	if err != nil {
		ctrl.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, playlist)
}

func (ctrl *Controller) UpdatePlaylist(w http.ResponseWriter, r *http.Request) {
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

	playlistID, err := getIDFromRequest(r, "id")
	if err != nil {
		ctrl.handleError(w, err)
		return
	}

	var input dto.UpdatePlaylistRequest
	if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
		ctrl.handleError(w, errors.Join(errs.ErrInvalidRequestBody, err))
		return
	}

	playlist, err := ctrl.service.UpdatePlaylist(userID, userRole, domain.Playlist{
		ID: playlistID,
		Name: input.Name,
		Description: input.Description,
	})
	if err != nil {
		ctrl.handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, playlist)
}


func (ctrl *Controller) DeletePlaylist(w http.ResponseWriter, r *http.Request) {
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

	playlistID, err := getIDFromRequest(r, "id")
	if err != nil {
		ctrl.handleError(w, err)
		return
	}
	if err = ctrl.service.DeletePlaylist(userID, userRole, playlistID); err != nil {
		ctrl.handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, CommonResponse{Message: "Playlist deleted successfully"})
}

func (ctrl *Controller) AddVideoToPlaylist(w http.ResponseWriter, r *http.Request) {
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

	playlistID, err := getIDFromRequest(r, "id")
	if err != nil {
		ctrl.handleError(w, err)
		return
	}

	var input dto.AddVideoToPlaylistRequest
	if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
		ctrl.handleError(w, errors.Join(errs.ErrInvalidRequestBody, err))
		return
	}
	if err = ctrl.service.AddVideoToPlaylist(userID, userRole, playlistID, input.VideoID); err != nil {
		ctrl.handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, CommonResponse{Message: "Video added to playlist successfully"})
}


func (ctrl *Controller) RemoveVideoFromPlaylist(w http.ResponseWriter, r *http.Request) {
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

	playlistID, err := getIDFromRequest(r, "id")
	if err != nil {
		ctrl.handleError(w, err)
		return
	}

	videoID, err := getIDFromRequest(r, "video_id")
	if err != nil {
		ctrl.handleError(w, err)
		return
	}
	if err = ctrl.service.RemoveVideoFromPlaylist(userID, userRole, playlistID, videoID); err != nil {
		ctrl.handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, CommonResponse{Message: "Video removed from playlist successfully"})
}