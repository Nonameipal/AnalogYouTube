package httpdelivery

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Nonameipal/AnalogYouTube/internal/domain"
	"github.com/Nonameipal/AnalogYouTube/internal/errs"
	appLogger "github.com/Nonameipal/AnalogYouTube/internal/logger"
	"github.com/gorilla/mux"
)

func (h *Handler) CreateVideo(w http.ResponseWriter, r *http.Request) {
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

	if err := r.ParseMultipartForm(600 << 20); err != nil {
		h.handleError(w, errs.ErrInvalidRequestBody)
		return
	}

	title := r.FormValue("title")
	description := r.FormValue("description")
	categoryName := r.FormValue("category_name")

	videoURL, err := h.saveMultipartFile(r, "video", "videos/originals", true)
	if err != nil {
		h.handleError(w, err)
		return
	}

	thumbnailURL, err := h.saveMultipartFile(r, "thumbnail", "thumbnails", true)
	if err != nil {
		h.cleanupFailedVideoUpload(0, "", 0, "", videoURL)
		h.handleError(w, err)
		return
	}

	var categoryNamePointer *string
	if categoryName != "" {
		categoryNamePointer = &categoryName
	}

	video, err := h.service.CreateVideo(userID, domain.Video{
		Title:        title,
		Description:  description,
		VideoURL:     videoURL,
		ThumbnailURL: thumbnailURL,
		CategoryName: categoryNamePointer,
	})
	if err != nil {
		h.cleanupFailedVideoUpload(0, "", 0, "", videoURL, thumbnailURL)
		h.handleError(w, err)
		return
	}

	inputPath := h.storage.URLPath(video.VideoURL)
	outputDir, outputURLPrefix := h.storage.VideoQualitiesPaths(video.ID)

	qualities, err := h.service.GenerateVideoQualities(userID, userRole, video.ID, inputPath, outputDir, outputURLPrefix)
	if err != nil {
		h.cleanupFailedVideoUpload(userID, userRole, video.ID, outputDir, video.VideoURL, video.ThumbnailURL)
		h.handleError(w, err)
		return
	}

	video.Qualities = qualities

	writeJSON(w, http.StatusCreated, video)
}

func (h *Handler) GetAllVideos(w http.ResponseWriter, r *http.Request) {
	videos, err := h.service.GetAllVideos()
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, videos)
}

func (h *Handler) SearchVideosByTitle(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("title")

	videos, err := h.service.SearchVideosByTitle(title)
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, videos)
}

func (h *Handler) GetRecommendedVideos(w http.ResponseWriter, r *http.Request) {
	videos, err := h.service.GetRecommendedVideos()
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, videos)
}

func (h *Handler) GetVideoByID(w http.ResponseWriter, r *http.Request) {
	videoID, err := getIDFromRequest(r, "id")
	if err != nil {
		h.handleError(w, err)
		return
	}

	if err = h.service.IncrementVideoViews(videoID); err != nil {
		h.handleError(w, err)
		return
	}

	video, err := h.service.GetVideoByID(videoID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, video)
}

func (h *Handler) UpdateVideo(w http.ResponseWriter, r *http.Request) {
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

	videoID, err := getIDFromRequest(r, "id")
	if err != nil {
		h.handleError(w, err)
		return
	}

	if err = r.ParseMultipartForm(50 << 20); err != nil {
		h.handleError(w, errs.ErrInvalidRequestBody)
		return
	}

	title := r.FormValue("title")
	description := r.FormValue("description")
	categoryName := r.FormValue("category_name")

	thumbnailURL, err := h.saveMultipartFile(r, "thumbnail", "thumbnails", false)
	if err != nil {
		h.handleError(w, err)
		return
	}

	var categoryNamePointer *string
	if categoryName != "" {
		categoryNamePointer = &categoryName
	}

	video, err := h.service.UpdateVideo(userID, userRole, domain.Video{
		ID:           videoID,
		Title:        title,
		Description:  description,
		ThumbnailURL: thumbnailURL,
		CategoryName: categoryNamePointer,
	})
	if err != nil {
		if thumbnailURL != "" {
			h.cleanupFailedVideoUpload(0, "", 0, "", thumbnailURL)
		}
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, video)
}

func (h *Handler) DeleteVideo(w http.ResponseWriter, r *http.Request) {
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

	videoID, err := getIDFromRequest(r, "id")
	if err != nil {
		h.handleError(w, err)
		return
	}

	video, err := h.service.GetVideoByID(videoID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	inputPath := h.storage.URLPath(video.VideoURL)
	archivePath, archiveURL := h.storage.VideoArchivePath(video.ID)

	if err = h.service.DeleteVideoWithArchive(userID, userRole, videoID, inputPath, archivePath, archiveURL); err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, CommonResponse{Message: "Video deleted and archived successfully"})
}

func (h *Handler) AdminDeleteVideo(w http.ResponseWriter, r *http.Request) {
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

	videoID, err := getIDFromRequest(r, "id")
	if err != nil {
		h.handleError(w, err)
		return
	}

	if err = h.service.DeleteVideo(userID, userRole, videoID); err != nil {
		h.handleError(w, err)
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

func (h *Handler) saveMultipartFile(r *http.Request, fieldName string, folder string, required bool) (string, error) {
	file, header, err := r.FormFile(fieldName)
	if err != nil {
		if !required && errors.Is(err, http.ErrMissingFile) {
			return "", nil
		}
		return "", errs.ErrInvalidRequestBody
	}
	defer file.Close()

	return h.storage.SaveFile(file, header, folder)
}

func (h *Handler) cleanupFailedVideoUpload(userID int, userRole string, videoID int, outputDir string, fileURLs ...string) {
	logger := appLogger.GetLogger()

	if videoID > 0 {
		if err := h.service.DeleteVideo(userID, userRole, videoID); err != nil {
			logger.Error().Err(err).Int("video_id", videoID).Msg("failed to delete video after upload error")
		}
	}

	for _, fileURL := range fileURLs {
		if err := h.storage.RemoveFileByURL(fileURL); err != nil {
			logger.Error().Err(err).Str("file_url", fileURL).Msg("failed to remove uploaded file after upload error")
		}
	}

	if err := h.storage.RemoveDir(outputDir); err != nil {
		logger.Error().Err(err).Str("output_dir", outputDir).Msg("failed to remove video qualities after upload error")
	}
}

func (h *Handler) GetPlaybackSpeeds(w http.ResponseWriter, r *http.Request) {
	speeds := h.service.GetPlaybackSpeeds()

	writeJSON(w, http.StatusOK, map[string][]float64{
		"speeds": speeds,
	})
}
