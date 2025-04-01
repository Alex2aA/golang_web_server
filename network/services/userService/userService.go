package userService

import (
	"context"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"idk_web_server001/network"
	"idk_web_server001/network/services/tokenService"
	"idk_web_server001/structures"
	"log"
	"net/http"
)

func searchUser(username string) (*structures.User, error) {
	rows, err := network.Pool.Query(context.Background(), "SELECT * FROM users WHERE username = $1", username)
	var user structures.User
	if err != nil {
		log.Println(err.Error())
		return &user, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&user.Id, &user.Name, &user.Password)
		if err != nil {
			log.Println(err.Error())
			return &user, err
		}
	}

	return &user, nil

}

func Login(username, password string) *structures.JSONMessage {
	user, err := searchUser(username)
	if err != nil {
		return &structures.JSONMessage{Status: http.StatusInternalServerError, Message: "Something went wrong"}
	}
	if user.Id == "" {
		return &structures.JSONMessage{Status: http.StatusUnauthorized, Message: "User " + username + " not found"}
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		return &structures.JSONMessage{Status: http.StatusBadRequest, Message: "Wrong Password"}
	}

	refreshToken, token, err := tokenService.GenerateTokens(user.Id)

	return &structures.JSONMessage{Status: http.StatusOK, Message: "Login", Token: token, RefreshToken: refreshToken}
}

func Register(username string, password string) *structures.JSONMessage {
	user, err := searchUser(username)
	if err != nil {
		return &structures.JSONMessage{Status: http.StatusInternalServerError, Message: "Something went wrong"}
	}
	if user.Id != "" {
		return &structures.JSONMessage{Status: http.StatusConflict, Message: "User already exists"}
	}
	id, err := uuid.NewRandom()

	if err != nil {
		log.Println(err.Error())
		return &structures.JSONMessage{Status: http.StatusInternalServerError, Message: "Something went wrong"}
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	if err != nil {
		log.Println(err.Error())
		return &structures.JSONMessage{Status: http.StatusInternalServerError, Message: "Something went wrong"}
	}
	_, err = network.Pool.Exec(context.Background(), "INSERT INTO users (id, username, hash_password) VALUES ($1, $2, $3)", id.String(), username, hashPassword)
	if err != nil {
		log.Println(err.Error())
		return &structures.JSONMessage{Status: http.StatusInternalServerError, Message: "Something went wrong"}
	}

	refreshToken, token, err := tokenService.GenerateTokens(id.String())
	if err != nil {
		return &structures.JSONMessage{Status: http.StatusInternalServerError, Message: "Create token error"}
	}

	return &structures.JSONMessage{Status: http.StatusCreated, Message: "User created", Token: token, RefreshToken: refreshToken}
}
