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
	rows, err := network.Pool.Query(context.Background(), "SELECT * FROM postes where user_id = $1", userId)

	var postes []structures.Post

	if err != nil {
		log.Println(err.Error())
		return postes, err
	}
	for rows.Next() {
		var newPost = structures.Post{}
		err = rows.Scan(&newPost.Id, &newPost.Name, &newPost.Description, &newPost.Text, &newPost.UserId)
		if err != nil {
			log.Println(err.Error())
			return postes, err
		}
		postes = append(postes, newPost)
	}

	return postes, nil
}

func getComments(postId string) ([]struct {
	userId string
	text   string
}, error) {
	rows, err := network.Pool.Query(context.Background(), "SELECT user_id, text FROM comments where post_id = $1", postId)
	var comments []struct {
		userId string
		text   string
	}
	if err != nil {
		return comments, err
	}
	for rows.Next() {
		var comment = struct {
			userId string
			text   string
		}{}
		err = rows.Scan(&comment.userId, &comment.text)
		if err != nil {
			return comments, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

func GetPost(userId string, postId string) (structures.Post, error) {
	row := network.Pool.QueryRow(context.Background(), "SELECT * FROM postes WHERE post_id = $1", postId)
	var post structures.Post
	println(row)
	return post, nil
}

func SendComment(userId string, postId string, text string) (*structures.JSONMessage, error) {
	_, err := network.Pool.Exec(context.Background(), "INSERT INTO comments_postes (post_id, user_id, text) VALUES ($1, $2, $3)", postId, userId, text)
	if err != nil {
		return &structures.JSONMessage{Status: http.StatusBadRequest, Message: "Bad Request"}, err
	}
	return &structures.JSONMessage{Status: http.StatusCreated, Message: "Comment created"}, nil
}

func CreatePost(post *structures.Post, userId string) (*structures.JSONMessage, error) {
	id, err := uuid.NewRandom()
	message := structures.JSONMessage{}
	if err != nil {
		message.Status = http.StatusInternalServerError
		message.Message = err.Error()
		return &message, err
	}
	_, err = network.Pool.Exec(context.Background(), "INSERT INTO postes (id,name,description,text, user_id) VALUES($1,$2,$3,$4,$5)", id.String(), post.Name, post.Description, post.Text, userId)
	if err != nil {
		message.Status = http.StatusInternalServerError
		message.Message = err.Error()
		return &message, err
	}
	message.Status = http.StatusCreated
	message.Message = "Post created"
	return &message, nil
}
