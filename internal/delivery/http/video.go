package httpdelivery

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/Nonameipal/AnalogYouTube/internal/delivery/http/dto"
	"github.com/Nonameipal/AnalogYouTube/internal/domain"
	"github.com/Nonameipal/AnalogYouTube/internal/errs"
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
		Title: input.Title,
		Description: input.Description,
		VideoURL: input.VideoURL,
		ThumbnailURL: input.ThumbnailURL,
		CategoryName: input.CategoryName,
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

func (ctrl *Controller) SearchVideosByTitle(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("title")

	videos, err := ctrl.service.SearchVideosByTitle(title)
	if err != nil {
		ctrl.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, videos)
}

func (ctrl *Controller) GetRecommendedVideos(w http.ResponseWriter, r *http.Request) {
	videos, err := ctrl.service.GetRecommendedVideos()
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

	if err = ctrl.service.IncrementVideoViews(videoID); err != nil {
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
		ID: videoID,
		Title: input.Title,
		Description: input.Description,
		VideoURL: input.VideoURL,
		ThumbnailURL: input.ThumbnailURL,
		CategoryName: input.CategoryName,
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

func (ctrl *Controller) UploadVideoThumbnail(w http.ResponseWriter, r *http.Request) {
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

	if err = r.ParseMultipartForm(10 << 20); err != nil {
		ctrl.handleError(w, errs.ErrInvalidRequestBody)
		return
	}

	file, header, err := r.FormFile("thumbnail")
	if err != nil {
		ctrl.handleError(w, errs.ErrInvalidRequestBody)
		return
	}
	defer file.Close()

	thumbnailURL, err := ctrl.storage.SaveFile(file, header, "thumbnails")
	if err != nil {
		ctrl.handleError(w, err)
		return
	}

	video, err := ctrl.service.UpdateVideoThumbnail(userID, userRole, videoID, thumbnailURL)
	if err != nil {
		ctrl.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, video)
}

func (ctrl *Controller) UploadVideoFile(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		ctrl.handleError(w, errs.ErrInvalidToken)
		return
	}

	if err := r.ParseMultipartForm(500 << 20); err != nil {
		ctrl.handleError(w, errs.ErrInvalidRequestBody)
		return
	}

	title := r.FormValue("title")
	description := r.FormValue("description")
	categoryName := r.FormValue("category_name")

	file, header, err := r.FormFile("video")
	if err != nil {
		ctrl.handleError(w, errs.ErrInvalidRequestBody)
		return
	}
	defer file.Close()

	videoURL, err := ctrl.storage.SaveFile(file, header, "videos/originals")
	if err != nil {
		ctrl.handleError(w, err)
		return
	}

	var categoryNamePointer *string
	if categoryName != "" {
		categoryNamePointer = &categoryName
	}

	video, err := ctrl.service.CreateVideo(userID, domain.Video{
	Title: title,
	Description: description,
	VideoURL: videoURL,
	CategoryName: categoryNamePointer,
})
if err != nil {
	ctrl.handleError(w, err)
	return
}

userRole, ok := getUserRoleFromContext(r)
if !ok {
	ctrl.handleError(w, errs.ErrInvalidToken)
	return
}

inputPath := ctrl.storage.URLPath(video.VideoURL)
outputDir, outputURLPrefix := ctrl.storage.VideoQualitiesPaths(video.ID)

qualities, err := ctrl.service.GenerateVideoQualities(userID, userRole, video.ID, inputPath, outputDir, outputURLPrefix)
if err != nil {
	ctrl.handleError(w, err)
	return
}

video.Qualities = qualities

writeJSON(w, http.StatusCreated, video)
}