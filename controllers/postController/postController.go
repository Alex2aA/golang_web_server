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

func GetPost(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("userId").(string)
	if !ok {
		httpHandlers.SendJSONMessage(w, structures.JSONMessage{Status: http.StatusInternalServerError, Message: "userId required"})
		return
	}
	var postId struct {
		PostId string `json:"post_id"`
	}
	err := json.NewDecoder(r.Body).Decode(&postId.PostId)
	defer r.Body.Close()
	if err != nil {
		httpHandlers.SendJSONMessage(w, structures.JSONMessage{Status: http.StatusBadRequest, Message: "Bad Request"})
		return
	}
	post, err := postService.GetPost(userId, postId.PostId)
	if err != nil {
		httpHandlers.SendJSONMessage(w, structures.JSONMessage{Status: http.StatusNotFound, Message: err.Error()})
		return
	}
	httpHandlers.SendJSONMessage(w, structures.JSONMessage{Status: http.StatusOK, Post: post})
}

func GetMyPostes(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("userId").(string)
	if !ok {
		httpHandlers.SendJSONMessage(w, structures.JSONMessage{Status: http.StatusInternalServerError, Message: "userId required"})
		return
	}
	result, err := postService.GetMyPostes(userId)
	if err != nil {
		httpHandlers.SendJSONMessage(w, structures.JSONMessage{Status: http.StatusInternalServerError, Message: err.Error()})
		return
	}
	httpHandlers.SendJSONMessage(w, structures.JSONMessage{Status: http.StatusOK, Message: "Get postes", Postes: result})
}

func SendComment(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("userId").(string)
	if !ok {
		httpHandlers.SendJSONMessage(w, structures.JSONMessage{Status: http.StatusInternalServerError, Message: "userId required"})
	}
	var comment struct {
		Text   string `json:"text"`
		PostId string `json:"post_id"`
	}
	err := json.NewDecoder(r.Body).Decode(&comment)
	defer r.Body.Close()
	if err != nil {
		httpHandlers.SendJSONMessage(w, structures.JSONMessage{Status: http.StatusBadRequest, Message: "Bad Request"})
		return
	}
	result, err := postService.SendComment(userId, comment.PostId, comment.Text)
	if err != nil {
		httpHandlers.SendJSONMessage(w, structures.JSONMessage{Status: result.Status, Message: result.Message})
		return
	}
	httpHandlers.SendJSONMessage(w, structures.JSONMessage{Status: result.Status, Message: result.Message})
}
