package httpdelivery

import (
	"fmt"
	"net/http"

	"github.com/Nonameipal/AnalogYouTube/internal/configs"
	"github.com/Nonameipal/AnalogYouTube/internal/logger"
	"github.com/gorilla/mux"
)

func (h *Handler) InitRoutes() error {
	r := mux.NewRouter()

	r.HandleFunc("/ping", h.ping).Methods(http.MethodGet)
	r.HandleFunc("/swagger", h.swaggerUI)
	r.PathPrefix("/swagger/").HandlerFunc(h.swaggerUI)

	auth := r.PathPrefix("/auth").Subrouter()
	{
		auth.HandleFunc("/sign-up", h.SignUp).Methods(http.MethodPost)
		auth.HandleFunc("/sign-in", h.SignIn).Methods(http.MethodPost)
		auth.HandleFunc("/refresh", h.RefreshTokenPair).Methods(http.MethodGet)
	}

	api := r.PathPrefix("/api").Subrouter()
	{
		api.HandleFunc("/register", h.SignUp).Methods(http.MethodPost)
		api.HandleFunc("/login", h.SignIn).Methods(http.MethodPost)
		api.HandleFunc("/videos", h.GetRecommendedVideos).Methods(http.MethodGet)
		api.HandleFunc("/videos/search", h.SearchVideosByTitle).Methods(http.MethodGet)

		api.HandleFunc("/videos/playback-speeds", h.GetPlaybackSpeeds).Methods(http.MethodGet)
		api.HandleFunc("/videos/{id}", h.GetVideoByID).Methods(http.MethodGet)
		api.HandleFunc("/videos/{id}/likes/count", h.GetVideoLikesCount).Methods(http.MethodGet)
		api.HandleFunc("/videos/{id}/comments", h.GetVideoComments).Methods(http.MethodGet)

		api.HandleFunc("/categories", h.GetAllCategories).Methods(http.MethodGet)
		api.HandleFunc("/categories/{name}", h.GetCategoryByName).Methods(http.MethodGet)

		api.HandleFunc("/users/{id}", h.GetUserProfile).Methods(http.MethodGet)
		api.HandleFunc("/users/{id}/videos", h.GetUserVideos).Methods(http.MethodGet)
		api.HandleFunc("/users/{id}/donations", h.GetUserDonations).Methods(http.MethodGet)
		api.HandleFunc("/users/{id}/subscribers/count", h.GetSubscribersCount).Methods(http.MethodGet)
		api.HandleFunc("/users/{id}/subscriptions/count", h.GetSubscriptionsCount).Methods(http.MethodGet)

		api.HandleFunc("/users/{id}/playlists", h.GetUserPlaylists).Methods(http.MethodGet)
		api.HandleFunc("/playlists/{id}", h.GetPlaylistByID).Methods(http.MethodGet)
	}

	auth = api.PathPrefix("").Subrouter()
	auth.Use(h.checkUserAuthentication)
	{
		auth.HandleFunc("/me", h.Me).Methods(http.MethodGet)
		auth.HandleFunc("/me", h.UpdateMe).Methods(http.MethodPut)

		auth.HandleFunc("/videos", h.CreateVideo).Methods(http.MethodPost)
		auth.HandleFunc("/videos/{id}", h.UpdateVideo).Methods(http.MethodPut)
		auth.HandleFunc("/videos/{id}", h.DeleteVideo).Methods(http.MethodDelete)

		auth.HandleFunc("/videos/{id}/like", h.LikeVideo).Methods(http.MethodPost)
		auth.HandleFunc("/videos/{id}/like", h.UnlikeVideo).Methods(http.MethodDelete)
		auth.HandleFunc("/videos/{id}/liked", h.IsVideoLikedByMe).Methods(http.MethodGet)

		auth.HandleFunc("/users/{id}/subscribe", h.SubscribeToUser).Methods(http.MethodPost)
		auth.HandleFunc("/users/{id}/subscribe", h.UnsubscribeFromUser).Methods(http.MethodDelete)
		auth.HandleFunc("/users/{id}/subscribed", h.IsSubscribedByMe).Methods(http.MethodGet)

		auth.HandleFunc("/videos/{id}/comments", h.CreateComment).Methods(http.MethodPost)
		auth.HandleFunc("/comments/{id}", h.UpdateComment).Methods(http.MethodPut)
		auth.HandleFunc("/comments/{id}", h.DeleteComment).Methods(http.MethodDelete)

		auth.HandleFunc("/donations", h.CreateDonation).Methods(http.MethodPost)
		auth.HandleFunc("/donations/sent", h.GetSentDonations).Methods(http.MethodGet)
		auth.HandleFunc("/donations/received", h.GetReceivedDonations).Methods(http.MethodGet)

		auth.HandleFunc("/chats", h.SendChatRequest).Methods(http.MethodPost)
		auth.HandleFunc("/chats", h.GetMyChats).Methods(http.MethodGet)
		auth.HandleFunc("/chats/requests/incoming", h.GetIncomingChatRequests).Methods(http.MethodGet)
		auth.HandleFunc("/chats/requests/outgoing", h.GetOutgoingChatRequests).Methods(http.MethodGet)
		auth.HandleFunc("/chats/requests/{id}/accept", h.AcceptChatRequest).Methods(http.MethodPost)
		auth.HandleFunc("/chats/requests/{id}/reject", h.RejectChatRequest).Methods(http.MethodPost)
		auth.HandleFunc("/chats/{id}/messages", h.GetChatMessages).Methods(http.MethodGet)
		auth.HandleFunc("/chats/{id}/messages", h.CreateChatMessage).Methods(http.MethodPost)

		auth.HandleFunc("/playlists", h.CreatePlaylist).Methods(http.MethodPost)
		auth.HandleFunc("/playlists/{id}", h.UpdatePlaylist).Methods(http.MethodPut)
		auth.HandleFunc("/playlists/{id}", h.DeletePlaylist).Methods(http.MethodDelete)
		auth.HandleFunc("/playlists/{id}/videos", h.AddVideoToPlaylist).Methods(http.MethodPost)
		auth.HandleFunc("/playlists/{id}/videos/{video_id}", h.RemoveVideoFromPlaylist).Methods(http.MethodDelete)
	}
	admin := api.PathPrefix("").Subrouter()
	admin.Use(h.checkUserAuthentication)
	admin.Use(h.checkIsAdmin)
	{
		admin.HandleFunc("/categories", h.CreateCategory).Methods(http.MethodPost)
		admin.HandleFunc("/categories/{id}", h.UpdateCategory).Methods(http.MethodPut)
		admin.HandleFunc("/categories/{id}", h.DeleteCategory).Methods(http.MethodDelete)
		admin.HandleFunc("/videos/allvideos", h.GetAllVideos).Methods(http.MethodGet)
		admin.HandleFunc("/admin/videos/{id}", h.AdminDeleteVideo).Methods(http.MethodDelete)
	}

	uploadsHandler := http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads")))
	r.PathPrefix("/uploads/").Handler(uploadsHandler)

	r.HandleFunc("/ws/chats/{id}", h.ChatWebSocket).Methods(http.MethodGet)

	addr := fmt.Sprintf(":%s", configs.AppSettings.AppParams.PortRun)
	logger.GetLogger().Info().Str("addr", addr).Msg("server started")

	return http.ListenAndServe(addr, r)
}

// ping godoc
// @Summary Проверка сервера
// @Description Быстрый ответ, что backend жив.
// @Tags System
// @Produce json
// @Success 200 {object} CommonResponse
// @Router /ping [get]
func (h *Handler) ping(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, CommonResponse{Message: "Server is running"})
}
