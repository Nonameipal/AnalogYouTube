package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/Nonameipal/AnalogYouTube/internal/errs"
	"github.com/Nonameipal/AnalogYouTube/internal/models/domain"
	"github.com/Nonameipal/AnalogYouTube/internal/models/dto"
	"github.com/gorilla/mux"
)

func (ctrl *Controller) CreateVideo(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		ctrl.handleError(w, errs.ErrInvalidToken)
		return
	}

	var input dto.CreateVideoRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		ctrl.handleError(w, errors.Join(errs.ErrInvalidRequestBody, err))
		return
	}

	video, err := ctrl.service.CreateVideo(userID, domain.Video{
		Title:        input.Title,
		Description:  input.Description,
		VideoURL:     input.VideoURL,
		ThumbnailURL: input.ThumbnailURL,
	})
	if err != nil {
		ctrl.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, video)
}

func (ctrl *Controller) GetAllVideos(w http.ResponseWriter, r *http.Request) {
	videos, err := ctrl.service.GetAllVideos()
	if err != nil {
		ctrl.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, videos)
}

func (ctrl *Controller) GetVideoByID(w http.ResponseWriter, r *http.Request) {
	videoID, err := getIDFromRequest(r, "id")
	if err != nil {
		ctrl.handleError(w, err)
		return
	}

	video, err := ctrl.service.GetVideoByID(videoID)
	if err != nil {
		ctrl.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, video)
}

func (ctrl *Controller) UpdateVideo(w http.ResponseWriter, r *http.Request) {
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

	videoID, err := getIDFromRequest(r, "id")
	if err != nil {
		ctrl.handleError(w, err)
		return
	}

	var input dto.UpdateVideoRequest
	if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
		ctrl.handleError(w, errors.Join(errs.ErrInvalidRequestBody, err))
		return
	}

	video, err := ctrl.service.UpdateVideo(userID, userRole, domain.Video{
		ID:           videoID,
		Title:        input.Title,
		Description:  input.Description,
		VideoURL:     input.VideoURL,
		ThumbnailURL: input.ThumbnailURL,
	})
	if err != nil {
		ctrl.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, video)
}

func (ctrl *Controller) DeleteVideo(w http.ResponseWriter, r *http.Request) {
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

	videoID, err := getIDFromRequest(r, "id")
	if err != nil {
		ctrl.handleError(w, err)
		return
	}

	if err = ctrl.service.DeleteVideo(userID, userRole, videoID); err != nil {
		ctrl.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, CommonResponse{Message: "Video deleted successfully"})
}

func getIDFromRequest(r *http.Request, key string) (int, error) {
	id, err := strconv.Atoi(mux.Vars(r)[key])
	if err != nil || id <= 0 {
		return 0, errs.ErrInvalidFieldValue
	}

	return id, nil
}
