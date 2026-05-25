package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Nonameipal/AnalogYouTube/internal/configs"
	"github.com/gorilla/mux"
)

func (ctrl *Controller) InitRoutes() error {
	r := mux.NewRouter()

	// понг
	r.HandleFunc("/ping", ctrl.ping).Methods(http.MethodGet)

	apiV1 := r.PathPrefix("/api/v1").Subrouter()


	// Публичные роуты
	apiV1.HandleFunc("/register", ctrl.Register).Methods(http.MethodPost)
	apiV1.HandleFunc("/login", ctrl.Login).Methods(http.MethodPost)

	// Получение профиля пользователя
	apiV1.HandleFunc("/users/{name}", ctrl.GetUserByName).Methods(http.MethodGet)


	// Публичные роуты для видео
	apiV1.HandleFunc("/videos", ctrl.GetAllVideos).Methods(http.MethodGet)
	apiV1.HandleFunc("/videos/{id}", ctrl.GetVideoByid).Methods(http.MethodGet)

	// Чат во время трансляции
	r.HandleFunc("/ws/chat/{videoId}", ctrl.VideoChatWebSocket).Methods(http.MethodGet)

	// Приватный чат с другом
	r.HandleFunc("/ws/chat/private/{chatId}", ctrl.PrivateChatWebSocket).Methods(http.MethodGet)



	authG := apiV1.PathPrefix("").Subrouter()
	authG.Use(ctrl.checkUserAuthentication)
	{
		// Пользователи
		authG.HandleFunc("/users/{id}", ctrl.UpdateUser).Methods(http.MethodPut)

		// Видео
		authG.HandleFunc("/videos", ctrl.UploadVideo).Methods(http.MethodPost)
		authG.HandleFunc("/videos/{id}", ctrl.UpdateVideo).Methods(http.MethodPut)
		authG.HandleFunc("/videos/{id}", ctrl.DeleteVideo).Methods(http.MethodDelete)

		// Взаимодействие
		authG.HandleFunc("/videos/{id}/like", ctrl.LikeVideo).Methods(http.MethodPost)
		authG.HandleFunc("/videos/{id}/like", ctrl.UnlikeVideo).Methods(http.MethodDelete)

		authG.HandleFunc("/users/{id}/subscribe", ctrl.SubscribeToUser).Methods(http.MethodPost)
		authG.HandleFunc("/users/{id}/subscribe", ctrl.UnsubscribeFromUser).Methods(http.MethodDelete)

		// Донаты
		authG.HandleFunc("/donations", ctrl.CreateDonation).Methods(http.MethodPost)
		authG.HandleFunc("/users/{id}/donations", ctrl.GetUserDonations).Methods(http.MethodGet)

		// Друзья
		authG.HandleFunc("/users/{id}/friends/request", ctrl.SendFriendRequest).Methods(http.MethodPost)
		authG.HandleFunc("/users/{id}/friends/accept", ctrl.AcceptFriendRequest).Methods(http.MethodPost)
		authG.HandleFunc("/users/{id}/friends/reject", ctrl.RejectOrDeleteFriend).Methods(http.MethodDelete)

		// Чаты
		authG.HandleFunc("/chats", ctrl.GetActiveChats).Methods(http.MethodGet)
	}

	adminG := apiV1.PathPrefix("/admin").Subrouter()
	adminG.Use(ctrl.checkUserAuthentication)
	adminG.Use(ctrl.checkIsAdmin)
	{
		adminG.HandleFunc("/users", ctrl.AdminGetAllUsers).Methods(http.MethodGet)
		adminG.HandleFunc("/users/{id}", ctrl.AdminDeleteUser).Methods(http.MethodDelete)

		adminG.HandleFunc("/videos", ctrl.AdminGetAllVideos).Methods(http.MethodGet)
		adminG.HandleFunc("/videos/{id}", ctrl.AdminDeleteVideo).Methods(http.MethodDelete)

		adminG.HandleFunc("/donations", ctrl.AdminGetAllDonations).Methods(http.MethodGet)
	}

	addr := fmt.Sprintf(":%s", configs.AppSettings.AppParams.PortRun)

	err := http.ListenAndServe(addr, r)
	if err != nil {
		return err
	}

	return nil
}

func (ctrl *Controller) ping(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, CommonResponse{
		Message: "Server is running",
	})
}

func writeJSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}