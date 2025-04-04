package postController

import (
	"encoding/json"
	"golang_web_server/network/httpHandlers"
	"golang_web_server/network/services/postService"
	"golang_web_server/structures"
	"log"
	"net/http"
)

func CreatePost(w http.ResponseWriter, r *http.Request) {
	var post structures.Post
	err := json.NewDecoder(r.Body).Decode(&post)
	defer r.Body.Close()
	if err != nil {
		log.Println(err.Error())
		httpHandlers.SendJSONMessage(w, structures.JSONMessage{Status: http.StatusBadRequest, Message: "Bad Request"})
		return
	}
	userId, ok := r.Context().Value("userId").(string)
	if !ok {
		log.Println(err.Error())
		httpHandlers.SendJSONMessage(w, structures.JSONMessage{Status: http.StatusInternalServerError, Message: "userId required"})
		return
	}
	result, err := postService.CreatePost(&post, userId)
	if err != nil {
		httpHandlers.SendJSONMessage(w, structures.JSONMessage{Status: result.Status, Message: result.Message})
		return
	}

	httpHandlers.SendJSONMessage(w, structures.JSONMessage{Status: result.Status, Message: result.Message})
}
