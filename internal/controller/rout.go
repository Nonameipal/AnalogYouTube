package controller

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

	authG := r.PathPrefix("/auth").Subrouter()
	{
		authG.HandleFunc("/sign-up", ctrl.SignUp).Methods(http.MethodPost)
		authG.HandleFunc("/sign-in", ctrl.SignIn).Methods(http.MethodPost)
		authG.HandleFunc("/refresh", ctrl.RefreshTokenPair).Methods(http.MethodGet)
	}

	apiV1 := r.PathPrefix("/api/v1").Subrouter()
	{
		apiV1.HandleFunc("/register", ctrl.SignUp).Methods(http.MethodPost)
		apiV1.HandleFunc("/login", ctrl.SignIn).Methods(http.MethodPost)

		apiV1.HandleFunc("/videos", ctrl.GetAllVideos).Methods(http.MethodGet)
		apiV1.HandleFunc("/videos/{id}", ctrl.GetVideoByID).Methods(http.MethodGet)
	}

	authApiV1 := apiV1.PathPrefix("").Subrouter()
	authApiV1.Use(ctrl.checkUserAuthentication)
	{
		authApiV1.HandleFunc("/me", ctrl.Me).Methods(http.MethodGet)

		authApiV1.HandleFunc("/videos", ctrl.CreateVideo).Methods(http.MethodPost)
		authApiV1.HandleFunc("/videos/{id}", ctrl.UpdateVideo).Methods(http.MethodPut)
		authApiV1.HandleFunc("/videos/{id}", ctrl.DeleteVideo).Methods(http.MethodDelete)
	}

	addr := fmt.Sprintf(":%s", configs.AppSettings.AppParams.PortRun)
	appLogger.GetLogger().Info().Str("addr", addr).Msg("server started")

	return http.ListenAndServe(addr, r)
}

func (ctrl *Controller) ping(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, CommonResponse{Message: "Server is running"})
}
