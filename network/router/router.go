package router

import (
	"github.com/gorilla/mux"
	"golang_web_server/middlewares/auth"
	"golang_web_server/network/httpHandlers"
)

func SetupRouter() *mux.Router {
	r := mux.NewRouter()

	api := r.PathPrefix("/api").Subrouter()

	api.HandleFunc("/health_check", httpHandlers.HealthCheck).Methods("GET")
	api.HandleFunc("/users/register", httpHandlers.Register).Methods("POST")
	api.HandleFunc("/users/login", httpHandlers.Login).Methods("POST")

	apiUploads := api.PathPrefix("/upload").Subrouter()

	apiUploads.Use(auth.AuthMiddleware)
	apiUploads.HandleFunc("/file/{typeFile}", httpHandlers.Upload).Methods("POST")

	//http.Handle("/avatars/", http.StripPrefix("/avatars/", http.FileServer(http.Dir(os.Getenv("UPLOAD_AVATARS_PATH")))))

	return r
}
