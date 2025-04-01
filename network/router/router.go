package router

import (
	"github.com/gorilla/mux"
	"golang_web_server/middlewares/auth"
	"golang_web_server/network/httpHandlers"
	"golang_web_server/network/httpHandlers/userHandler"
)

func SetupRouter() *mux.Router {
	r := mux.NewRouter()

	api := r.PathPrefix("/api").Subrouter()

	api.HandleFunc("/health_check", httpHandlers.HealthCheck).Methods("GET")
	api.HandleFunc("/users/register", userHandler.Register).Methods("POST")
	api.HandleFunc("/users/login", userHandler.Login).Methods("POST")

	apiUploads := api.PathPrefix("/upload").Subrouter()

	apiUploads.Use(auth.AuthMiddleware)
	apiUploads.HandleFunc("/file/{typeFile}", userHandler.Upload).Methods("POST")

	//http.Handle("/avatars/", http.StripPrefix("/avatars/", http.FileServer(http.Dir(os.Getenv("UPLOAD_AVATARS_PATH")))))

	return r
}
