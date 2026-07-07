package httpdelivery

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/Nonameipal/AnalogYouTube/internal/domain"
	"github.com/Nonameipal/AnalogYouTube/internal/errs"
	appLogger "github.com/Nonameipal/AnalogYouTube/internal/logger"
	"github.com/gorilla/mux"
)

// CreateVideo godoc
// @Summary Создать видео
// @Description Загружает видеофайл, обложку и данные ролика.
// @Tags Videos
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param title formData string true "Название видео"
// @Param description formData string false "Описание видео"
// @Param category_name formData string false "Имя категории"
// @Param video formData file true "Файл видео"
// @Param thumbnail formData file true "Файл обложки"
// @Success 201 {object} domain.Video
// @Failure 400 {object} CommonError
// @Failure 401 {object} CommonError
// @Router /api/videos [post]
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
	tagNames := parseTagNames(r.FormValue("tags"))

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
		Tags:         tagsFromNames(tagNames),
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

// GetAllVideos godoc
// @Summary Все видео
// @Description Админский список всех видео.
// @Tags Admin,Videos
// @Produce json
// @Security BearerAuth
// @Success 200 {array} domain.Video
// @Failure 401 {object} CommonError
// @Failure 403 {object} CommonError
// @Router /api/videos/allvideos [get]
func (h *Handler) GetAllVideos(w http.ResponseWriter, r *http.Request) {
	videos, err := h.service.GetAllVideos()
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, videos)
}

// SearchVideosByTitle godoc
// @Summary Поиск видео
// @Description Ищет ролики по названию.
// @Tags Videos
// @Produce json
// @Param title query string true "Часть названия видео"
// @Success 200 {array} domain.Video
// @Failure 400 {object} CommonError
// @Router /api/videos/search [get]
func (h *Handler) SearchVideosByTitle(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("title")
	userID := h.optionalUserIDFromRequest(r)

	videos, err := h.service.SearchVideosByTitleForUser(userID, title)
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, videos)
}

// GetRecommendedVideos godoc
// @Summary Рекомендации
// @Description Видео для главной ленты.
// @Tags Videos
// @Produce json
// @Success 200 {array} domain.Video
// @Router /api/videos [get]
func (h *Handler) GetRecommendedVideos(w http.ResponseWriter, r *http.Request) {
	userID := h.optionalUserIDFromRequest(r)
	if userID == nil {
		videos, err := h.service.GetRecommendedVideos()
		if err != nil {
			h.handleError(w, err)
			return
		}

		writeJSON(w, http.StatusOK, videos)
		return
	}

	videos, err := h.service.GetPersonalizedRecommendedVideos(*userID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, videos)
}

func (h *Handler) GetPersonalizedRecommendedVideos(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		h.handleError(w, errs.ErrInvalidToken)
		return
	}

	videos, err := h.service.GetPersonalizedRecommendedVideos(userID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, videos)
}

// GetVideoByID godoc
// @Summary Видео по ID
// @Description Данные видео и новый просмотр.
// @Tags Videos
// @Produce json
// @Param id path int true "ID видео"
// @Success 200 {object} domain.Video
// @Failure 400 {object} CommonError
// @Failure 404 {object} CommonError
// @Router /api/videos/{id} [get]
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

// UpdateVideo godoc
// @Summary Обновить видео
// @Description Меняет данные ролика без замены видеофайла.
// @Tags Videos
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID видео"
// @Param title formData string true "Название видео"
// @Param description formData string false "Описание видео"
// @Param category_name formData string false "Имя категории"
// @Param thumbnail formData file false "Новая обложка"
// @Success 200 {object} domain.Video
// @Failure 400 {object} CommonError
// @Failure 401 {object} CommonError
// @Failure 403 {object} CommonError
// @Router /api/videos/{id} [put]
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
	tagNames, tagsProvided := tagNamesFromMultipartForm(r, "tags")

	thumbnailURL, err := h.saveMultipartFile(r, "thumbnail", "thumbnails", false)
	if err != nil {
		h.handleError(w, err)
		return
	}

	var categoryNamePointer *string
	if categoryName != "" {
		categoryNamePointer = &categoryName
	}

	video := domain.Video{
		ID:           videoID,
		Title:        title,
		Description:  description,
		ThumbnailURL: thumbnailURL,
		CategoryName: categoryNamePointer,
	}
	if tagsProvided {
		video.Tags = tagsFromNames(tagNames)
	}
	video, err = h.service.UpdateVideo(userID, userRole, video)
	if err != nil {
		if thumbnailURL != "" {
			h.cleanupFailedVideoUpload(0, "", 0, "", thumbnailURL)
		}
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, video)
}

// DeleteVideo godoc
// @Summary Удалить видео
// @Description Скрывает видео и оставляет архивную копию.
// @Tags Videos
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID видео"
// @Success 200 {object} CommonResponse
// @Failure 401 {object} CommonError
// @Failure 403 {object} CommonError
// @Failure 404 {object} CommonError
// @Router /api/videos/{id} [delete]
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

// AdminDeleteVideo godoc
// @Summary Админское удаление
// @Description Полностью удаляет видео.
// @Tags Admin,Videos
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID видео"
// @Success 200 {object} CommonResponse
// @Failure 401 {object} CommonError
// @Failure 403 {object} CommonError
// @Router /api/admin/videos/{id} [delete]
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

// GetPlaybackSpeeds godoc
// @Summary Скорости видео
// @Description Доступные скорости плеера.
// @Tags Videos
// @Produce json
// @Success 200 {object} map[string][]float64
// @Router /api/videos/playback-speeds [get]
func (h *Handler) GetPlaybackSpeeds(w http.ResponseWriter, r *http.Request) {
	speeds := h.service.GetPlaybackSpeeds()

	writeJSON(w, http.StatusOK, map[string][]float64{
		"speeds": speeds,
	})
}

type watchProgressRequest struct {
	WatchedSeconds  int `json:"watched_seconds"`
	DurationSeconds int `json:"duration_seconds"`
}

type updateVideoTagsRequest struct {
	Tags []string `json:"tags"`
}

func (h *Handler) SaveVideoWatchProgress(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		h.handleError(w, errs.ErrInvalidToken)
		return
	}

	videoID, err := getIDFromRequest(r, "id")
	if err != nil {
		h.handleError(w, err)
		return
	}

	var request watchProgressRequest
	if err = json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.handleError(w, errs.ErrInvalidRequestBody)
		return
	}

	progress, err := h.service.SaveVideoWatchProgress(userID, videoID, request.WatchedSeconds, request.DurationSeconds)
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, progress)
}

func (h *Handler) GetAllTags(w http.ResponseWriter, r *http.Request) {
	tags, err := h.service.GetAllTags()
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, tags)
}

func (h *Handler) GetVideoTags(w http.ResponseWriter, r *http.Request) {
	videoID, err := getIDFromRequest(r, "id")
	if err != nil {
		h.handleError(w, err)
		return
	}

	tags, err := h.service.GetVideoTags(videoID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, tags)
}

func (h *Handler) UpdateVideoTags(w http.ResponseWriter, r *http.Request) {
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

	var request updateVideoTagsRequest
	if err = json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.handleError(w, errs.ErrInvalidRequestBody)
		return
	}

	tags, err := h.service.UpdateVideoTags(userID, userRole, videoID, request.Tags)
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, tags)
}

func tagNamesFromMultipartForm(r *http.Request, fieldName string) ([]string, bool) {
	if r.MultipartForm == nil {
		return nil, false
	}

	values, ok := r.MultipartForm.Value[fieldName]
	if !ok {
		return nil, false
	}

	return parseTagNames(strings.Join(values, ",")), true
}

func parseTagNames(value string) []string {
	if strings.TrimSpace(value) == "" {
		return []string{}
	}

	parts := strings.Split(value, ",")
	tags := make([]string, 0, len(parts))
	for _, part := range parts {
		if tag := strings.TrimSpace(part); tag != "" {
			tags = append(tags, tag)
		}
	}

	return tags
}

func tagsFromNames(names []string) []domain.Tag {
	if names == nil {
		return nil
	}

	tags := make([]domain.Tag, 0, len(names))
	for _, name := range names {
		tags = append(tags, domain.Tag{Name: name})
	}

	return tags
}
