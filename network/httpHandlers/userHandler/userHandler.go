package userHandler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"golang_web_server/network/httpHandlers"
	"golang_web_server/network/services/uploadAvatarService"
	"golang_web_server/network/services/userService"
	"golang_web_server/structures"
	"log"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {
	var user structures.User
	err := json.NewDecoder(r.Body).Decode(&user)
	defer r.Body.Close()
	if err != nil {
		log.Println(err.Error())
		httpHandlers.SendJSONMessage(w, structures.JSONMessage{Status: http.StatusBadRequest, Message: "Bad Request"})
		return
	}

	result := userService.Login(user.Name, user.Password)

	httpHandlers.SendJSONMessage(w, *result)
}

func Register(w http.ResponseWriter, r *http.Request) {
	var user structures.User
	err := json.NewDecoder(r.Body).Decode(&user)
	defer r.Body.Close()
	if err != nil {
		log.Println(err.Error())
		httpHandlers.SendJSONMessage(w, structures.JSONMessage{Status: http.StatusBadRequest, Message: "Bad Request"})
		return
	}
	result := userService.Register(user.Name, user.Password)
	httpHandlers.SendJSONMessage(w, *result)
}

func Upload(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	typeFile := vars["typeFile"]

	switch {
	case typeFile == "avatar":
		result := uploadAvatarService.InitUploadAvatarService(w, r)
		httpHandlers.SendJSONMessage(w, *result)
	case typeFile == "image":
		httpHandlers.SendJSONMessage(w, structures.JSONMessage{Status: http.StatusNotImplemented, Message: "Not Implemented"})
	default:
		httpHandlers.SendJSONMessage(w, structures.JSONMessage{Status: http.StatusNotImplemented, Message: "Something went wrong"})
	}
}
