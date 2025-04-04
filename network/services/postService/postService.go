package postService

import (
	"context"
	"github.com/google/uuid"
	"golang_web_server/network"
	"golang_web_server/structures"
	"log"
	"net/http"
)

func GetMyPostes(userId string) ([]structures.Post, error) {
	rows, err := network.Pool.Query(context.Background(), "SELECT * FROM users_postes WHERE user_id = $1", userId)
	postes := []structures.Post{}
	if err != nil {
		log.Println(err.Error())
		return postes, err
	}
	for rows.Next() {
		var newPost = structures.Post{}
		err = rows.Scan(&newPost.Id, &newPost.Name, &newPost.Description, &newPost.Text)
		if err != nil {
			log.Println(err.Error())
			return postes, err
		}
		//Доделать postes
	}
}

func CreatePost(post *structures.Post, userId string) (*structures.JSONMessage, error) {
	id, err := uuid.NewRandom()
	message := structures.JSONMessage{}
	if err != nil {
		message.Status = http.StatusInternalServerError
		message.Message = err.Error()
		return &message, err
	}
	_, err = network.Pool.Exec(context.Background(), "INSERT INTO postes (id,name,description,text) VALUES($1,$2,$3,$4)", id.String(), post.Name, post.Description, post.Text)
	if err != nil {
		message.Status = http.StatusInternalServerError
		message.Message = err.Error()
		return &message, err
	}
	_, err = network.Pool.Exec(context.Background(), "INSERT INTO users_postes (user_id, post_id) VALUES($1, $2)", userId, id.String())
	if err != nil {
		message.Status = http.StatusInternalServerError
		message.Message = err.Error()
		return &message, err
	}
	message.Status = http.StatusCreated
	message.Message = "Post created"
	return &message, nil
}
