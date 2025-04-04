package router

import (
	"github.com/gorilla/mux"
	"golang_web_server/controllers/postController"
	"golang_web_server/controllers/userController"
	"golang_web_server/middlewares/auth"
	"golang_web_server/network/httpHandlers"
)

func SetupRouter() *mux.Router {
	r := mux.NewRouter()

	api := r.PathPrefix("/api").Subrouter()

	api.HandleFunc("/health_check", httpHandlers.HealthCheck).Methods("GET")
	api.HandleFunc("/users/register", userController.Register).Methods("POST")
	api.HandleFunc("/users/login", userController.Login).Methods("POST")

	apiUploads := api.PathPrefix("/upload").Subrouter()

	apiUploads.Use(auth.AuthMiddleware)
	apiUploads.HandleFunc("/file/{typeFile}", userController.Upload).Methods("POST")

	post := api.PathPrefix("/post").Subrouter()
	post.Use(auth.AuthMiddleware)
	post.HandleFunc("/create_post", postController.CreatePost).Methods("POST")
	post.HandleFunc("/get_my_postes", postController.GetMyPostes).Methods("GET")

	//http.Handle("/avatars/", http.StripPrefix("/avatars/", http.FileServer(http.Dir(os.Getenv("UPLOAD_AVATARS_PATH")))))

	return r
}
