package httpdelivery

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Nonameipal/AnalogYouTube/internal/delivery/http/dto"
	"github.com/Nonameipal/AnalogYouTube/internal/domain"
	"github.com/Nonameipal/AnalogYouTube/internal/errs"
)

func (h *Handler) CreatePlaylist(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		h.handleError(w, errs.ErrInvalidToken)
		return
	}

	var input dto.CreatePlaylistRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.handleError(w, errors.Join(errs.ErrInvalidRequestBody, err))
		return
	}

	playlist, err := h.service.CreatePlaylist(userID, domain.Playlist{
		Name:        input.Name,
		Description: input.Description,
	})
	if err != nil {
		h.handleError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, playlist)
}

func (h *Handler) GetUserPlaylists(w http.ResponseWriter, r *http.Request) {
	userID, err := getIDFromRequest(r, "id")
	if err != nil {
		h.handleError(w, err)
		return
	}

	playlists, err := h.service.GetUserPlaylists(userID)
	if err != nil {
		h.handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, playlists)
}

func (h *Handler) GetPlaylistByID(w http.ResponseWriter, r *http.Request) {
	playlistID, err := getIDFromRequest(r, "id")
	if err != nil {
		h.handleError(w, err)
		return
	}

	playlist, err := h.service.GetPlaylistByID(playlistID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, playlist)
}

func (h *Handler) UpdatePlaylist(w http.ResponseWriter, r *http.Request) {
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

	playlistID, err := getIDFromRequest(r, "id")
	if err != nil {
		h.handleError(w, err)
		return
	}

	var input dto.UpdatePlaylistRequest
	if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.handleError(w, errors.Join(errs.ErrInvalidRequestBody, err))
		return
	}

	playlist, err := h.service.UpdatePlaylist(userID, userRole, domain.Playlist{
		ID:          playlistID,
		Name:        input.Name,
		Description: input.Description,
	})
	if err != nil {
		h.handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, playlist)
}

func (h *Handler) DeletePlaylist(w http.ResponseWriter, r *http.Request) {
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

	playlistID, err := getIDFromRequest(r, "id")
	if err != nil {
		h.handleError(w, err)
		return
	}
	if err = h.service.DeletePlaylist(userID, userRole, playlistID); err != nil {
		h.handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, CommonResponse{Message: "Playlist deleted successfully"})
}

func (h *Handler) AddVideoToPlaylist(w http.ResponseWriter, r *http.Request) {
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

	playlistID, err := getIDFromRequest(r, "id")
	if err != nil {
		h.handleError(w, err)
		return
	}

	var input dto.AddVideoToPlaylistRequest
	if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.handleError(w, errors.Join(errs.ErrInvalidRequestBody, err))
		return
	}
	if err = h.service.AddVideoToPlaylist(userID, userRole, playlistID, input.VideoID); err != nil {
		h.handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, CommonResponse{Message: "Video added to playlist successfully"})
}

func (h *Handler) RemoveVideoFromPlaylist(w http.ResponseWriter, r *http.Request) {
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

	playlistID, err := getIDFromRequest(r, "id")
	if err != nil {
		h.handleError(w, err)
		return
	}

	videoID, err := getIDFromRequest(r, "video_id")
	if err != nil {
		h.handleError(w, err)
		return
	}
	if err = h.service.RemoveVideoFromPlaylist(userID, userRole, playlistID, videoID); err != nil {
		h.handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, CommonResponse{Message: "Video removed from playlist successfully"})
}
