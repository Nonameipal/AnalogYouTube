package httpdelivery

import (
	"net/http"

	"github.com/Nonameipal/AnalogYouTube/internal/errs"
)

func (h *Handler) LikeVideo(w http.ResponseWriter, r *http.Request) {
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

	if err = h.service.LikeVideo(userID, videoID); err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, CommonResponse{Message: "Video liked successfully"})
}

func (h *Handler) UnlikeVideo(w http.ResponseWriter, r *http.Request) {
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

	if err = h.service.UnlikeVideo(userID, videoID); err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, CommonResponse{Message: "Video unliked successfully"})
}

func (h *Handler) GetVideoLikesCount(w http.ResponseWriter, r *http.Request) {
	videoID, err := getIDFromRequest(r, "id")
	if err != nil {
		h.handleError(w, err)
		return
	}

	count, err := h.service.GetVideoLikesCount(videoID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]int{"likes_count": count})
}

func (h *Handler) IsVideoLikedByMe(w http.ResponseWriter, r *http.Request) {
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

	liked, err := h.service.IsVideoLikedByUser(userID, videoID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]bool{"liked": liked})
}
