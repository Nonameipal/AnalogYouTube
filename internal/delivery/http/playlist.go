package httpdelivery

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Nonameipal/AnalogYouTube/internal/delivery/http/dto"
	"github.com/Nonameipal/AnalogYouTube/internal/domain"
	"github.com/Nonameipal/AnalogYouTube/internal/errs"
)

// CreatePlaylist godoc
// @Summary Создать плейлист
// @Description Новый публичный плейлист.
// @Tags Playlists
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body dto.CreatePlaylistRequest true "Данные плейлиста"
// @Success 201 {object} domain.Playlist
// @Failure 400 {object} CommonError
// @Failure 401 {object} CommonError
// @Router /api/playlists [post]
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

// GetUserPlaylists godoc
// @Summary Плейлисты пользователя
// @Description Публичные плейлисты автора.
// @Tags Playlists
// @Produce json
// @Param id path int true "ID пользователя"
// @Success 200 {array} domain.Playlist
// @Failure 400 {object} CommonError
// @Router /api/users/{id}/playlists [get]
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

// GetPlaylistByID godoc
// @Summary Плейлист по ID
// @Description Плейлист и его видео.
// @Tags Playlists
// @Produce json
// @Param id path int true "ID плейлиста"
// @Success 200 {object} domain.Playlist
// @Failure 404 {object} CommonError
// @Router /api/playlists/{id} [get]
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

// UpdatePlaylist godoc
// @Summary Обновить плейлист
// @Description Меняет название или описание.
// @Tags Playlists
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID плейлиста"
// @Param input body dto.UpdatePlaylistRequest true "Новые данные"
// @Success 200 {object} domain.Playlist
// @Failure 401 {object} CommonError
// @Failure 403 {object} CommonError
// @Router /api/playlists/{id} [put]
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

// DeletePlaylist godoc
// @Summary Удалить плейлист
// @Description Удаляет плейлист автора.
// @Tags Playlists
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID плейлиста"
// @Success 200 {object} CommonResponse
// @Failure 401 {object} CommonError
// @Failure 403 {object} CommonError
// @Router /api/playlists/{id} [delete]
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

// AddVideoToPlaylist godoc
// @Summary Добавить видео
// @Description Добавляет ролик в плейлист.
// @Tags Playlists
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID плейлиста"
// @Param input body dto.AddVideoToPlaylistRequest true "ID видео"
// @Success 200 {object} CommonResponse
// @Failure 401 {object} CommonError
// @Failure 403 {object} CommonError
// @Router /api/playlists/{id}/videos [post]
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

// RemoveVideoFromPlaylist godoc
// @Summary Убрать видео
// @Description Убирает ролик из плейлиста.
// @Tags Playlists
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID плейлиста"
// @Param video_id path int true "ID видео"
// @Success 200 {object} CommonResponse
// @Failure 401 {object} CommonError
// @Failure 403 {object} CommonError
// @Router /api/playlists/{id}/videos/{video_id} [delete]
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
