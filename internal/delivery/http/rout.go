package httpdelivery

import (
	"fmt"
	"net/http"

	"github.com/Nonameipal/AnalogYouTube/internal/configs"
	appLogger "github.com/Nonameipal/AnalogYouTube/internal/logger"
	"github.com/gorilla/mux"
)

func (ctrl *Controller) InitRoutes() error {
	r := mux.NewRouter()

	r.HandleFunc("/ping", ctrl.ping).Methods(http.MethodGet)

	auth := r.PathPrefix("/auth").Subrouter()
	{
		auth.HandleFunc("/sign-up", ctrl.SignUp).Methods(http.MethodPost)
		auth.HandleFunc("/sign-in", ctrl.SignIn).Methods(http.MethodPost)
		auth.HandleFunc("/refresh", ctrl.RefreshTokenPair).Methods(http.MethodGet)
	}

	api := r.PathPrefix("/api").Subrouter()
	{
		api.HandleFunc("/register", ctrl.SignUp).Methods(http.MethodPost)
		api.HandleFunc("/login", ctrl.SignIn).Methods(http.MethodPost)
		api.HandleFunc("/videos", ctrl.GetRecommendedVideos).Methods(http.MethodGet)
		api.HandleFunc("/videos/search", ctrl.SearchVideosByTitle).Methods(http.MethodGet)

		api.HandleFunc("/videos/{id}", ctrl.GetVideoByID).Methods(http.MethodGet)
		api.HandleFunc("/videos/{id}/likes/count", ctrl.GetVideoLikesCount).Methods(http.MethodGet)
		api.HandleFunc("/videos/{id}/comments", ctrl.GetVideoComments).Methods(http.MethodGet)

		api.HandleFunc("/categories", ctrl.GetAllCategories).Methods(http.MethodGet)
		api.HandleFunc("/categories/{name}", ctrl.GetCategoryByName).Methods(http.MethodGet)

		api.HandleFunc("/users/{id}", ctrl.GetUserProfile).Methods(http.MethodGet)
		api.HandleFunc("/users/{id}/videos", ctrl.GetUserVideos).Methods(http.MethodGet)
		api.HandleFunc("/users/{id}/donations", ctrl.GetUserDonations).Methods(http.MethodGet)
		api.HandleFunc("/users/{id}/subscribers/count", ctrl.GetSubscribersCount).Methods(http.MethodGet)
		api.HandleFunc("/users/{id}/subscriptions/count", ctrl.GetSubscriptionsCount).Methods(http.MethodGet)

		api.HandleFunc("/users/{id}/playlists", ctrl.GetUserPlaylists).Methods(http.MethodGet)
		api.HandleFunc("/playlists/{id}", ctrl.GetPlaylistByID).Methods(http.MethodGet)
	}

	auth = api.PathPrefix("").Subrouter()
	auth.Use(ctrl.checkUserAuthentication)
	{
		auth.HandleFunc("/me", ctrl.Me).Methods(http.MethodGet)
		auth.HandleFunc("/me", ctrl.UpdateMe).Methods(http.MethodPut)
		auth.HandleFunc("/me/avatar", ctrl.UploadMyAvatar).Methods(http.MethodPost)

		auth.HandleFunc("/videos", ctrl.CreateVideo).Methods(http.MethodPost)
		auth.HandleFunc("/videos/{id}", ctrl.UpdateVideo).Methods(http.MethodPut)
		auth.HandleFunc("/videos/{id}", ctrl.DeleteVideo).Methods(http.MethodDelete)
		auth.HandleFunc("/videos/{id}/thumbnail", ctrl.UploadVideoThumbnail).Methods(http.MethodPost)
		auth.HandleFunc("/videos/upload", ctrl.UploadVideoFile).Methods(http.MethodPost)

		auth.HandleFunc("/videos/{id}/like", ctrl.LikeVideo).Methods(http.MethodPost)
		auth.HandleFunc("/videos/{id}/like", ctrl.UnlikeVideo).Methods(http.MethodDelete)
		auth.HandleFunc("/videos/{id}/liked", ctrl.IsVideoLikedByMe).Methods(http.MethodGet)

		auth.HandleFunc("/users/{id}/subscribe", ctrl.SubscribeToUser).Methods(http.MethodPost)
		auth.HandleFunc("/users/{id}/subscribe", ctrl.UnsubscribeFromUser).Methods(http.MethodDelete)
		auth.HandleFunc("/users/{id}/subscribed", ctrl.IsSubscribedByMe).Methods(http.MethodGet)

		auth.HandleFunc("/videos/{id}/comments", ctrl.CreateComment).Methods(http.MethodPost)
		auth.HandleFunc("/comments/{id}", ctrl.UpdateComment).Methods(http.MethodPut)
		auth.HandleFunc("/comments/{id}", ctrl.DeleteComment).Methods(http.MethodDelete)

		auth.HandleFunc("/donations", ctrl.CreateDonation).Methods(http.MethodPost)
		auth.HandleFunc("/donations/sent", ctrl.GetSentDonations).Methods(http.MethodGet)
		auth.HandleFunc("/donations/received", ctrl.GetReceivedDonations).Methods(http.MethodGet)

		auth.HandleFunc("/chats", ctrl.CreateOrGetChat).Methods(http.MethodPost)
		auth.HandleFunc("/chats", ctrl.GetMyChats).Methods(http.MethodGet)
		auth.HandleFunc("/chats/{id}/messages", ctrl.GetChatMessages).Methods(http.MethodGet)
		auth.HandleFunc("/chats/{id}/messages", ctrl.CreateChatMessage).Methods(http.MethodPost)

		auth.HandleFunc("/playlists", ctrl.CreatePlaylist).Methods(http.MethodPost)
		auth.HandleFunc("/playlists/{id}", ctrl.UpdatePlaylist).Methods(http.MethodPut)
		auth.HandleFunc("/playlists/{id}", ctrl.DeletePlaylist).Methods(http.MethodDelete)
		auth.HandleFunc("/playlists/{id}/videos", ctrl.AddVideoToPlaylist).Methods(http.MethodPost)
		auth.HandleFunc("/playlists/{id}/videos/{video_id}", ctrl.RemoveVideoFromPlaylist).Methods(http.MethodDelete)
	}
	admin := api.PathPrefix("").Subrouter()
	admin.Use(ctrl.checkUserAuthentication)
	admin.Use(ctrl.checkIsAdmin)
	{
		admin.HandleFunc("/categories", ctrl.CreateCategory).Methods(http.MethodPost)
		admin.HandleFunc("/categories/{id}", ctrl.UpdateCategory).Methods(http.MethodPut)
		admin.HandleFunc("/categories/{id}", ctrl.DeleteCategory).Methods(http.MethodDelete)
		admin.HandleFunc("/videos/allvideos", ctrl.GetAllVideos).Methods(http.MethodGet)
	}

	uploadsHandler := http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads")))
	r.PathPrefix("/uploads/").Handler(uploadsHandler)

	r.HandleFunc("/ws/chats/{id}", ctrl.ChatWebSocket).Methods(http.MethodGet)

	addr := fmt.Sprintf(":%s", configs.AppSettings.AppParams.PortRun)
	appLogger.GetLogger().Info().Str("addr", addr).Msg("server started")

	return http.ListenAndServe(addr, r)
}

func (ctrl *Controller) ping(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, CommonResponse{Message: "Server is running"})
}
