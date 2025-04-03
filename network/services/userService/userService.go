package userService

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"golang_web_server/network"
	"golang_web_server/network/services/tokenService"
	"golang_web_server/structures"
	"log"
	"net/http"
)

func searchUserById(userId string) (*structures.User, error) {
	rows, err := network.Pool.Query(context.Background(), "SELECT * FROM users WHERE id = $1", userId)
	var user structures.User
	if err != nil {
		log.Println(err.Error())
		return &user, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&user.Id, &user.Name, &user.Password, &user.RefreshToken)
		if err != nil {
			log.Println(err.Error())
			return &user, err
		}
	}

	return &user, nil
}

func searchUserByName(username string) (*structures.User, error) {
	rows, err := network.Pool.Query(context.Background(), "SELECT * FROM users WHERE username = $1", username)
	var user structures.User
	if err != nil {
		log.Println(err.Error())
		return &user, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&user.Id, &user.Name, &user.Password, &user.RefreshToken)
		if err != nil {
			log.Println(err.Error())
			return &user, err
		}
	}

	return &user, nil

}

func Login(username, password string) *structures.JSONMessage {
	user, err := searchUserByName(username)
	if err != nil {
		return &structures.JSONMessage{Status: http.StatusInternalServerError, Message: err.Error()}
	}
	if user.Id == "" {
		return &structures.JSONMessage{Status: http.StatusUnauthorized, Message: "User " + username + " not found"}
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		return &structures.JSONMessage{Status: http.StatusBadRequest, Message: err.Error()}
	}

	tokenString, err := tokenService.GenerateAccessToken(user.Id)

	if err != nil {
		return &structures.JSONMessage{Status: http.StatusInternalServerError, Message: err.Error()}
	}

	refreshTokenString, err := tokenService.GenerateRefreshToken(user.Id)

	if err != nil {
		return &structures.JSONMessage{Status: http.StatusInternalServerError, Message: err.Error()}
	}

	_, err = network.Pool.Exec(context.Background(), "UPDATE users SET refresh_token = $1 WHERE id = $2", refreshTokenString, user.Id)

	if err != nil {
		return &structures.JSONMessage{Status: http.StatusInternalServerError, Message: err.Error()}
	}
	return &structures.JSONMessage{Status: http.StatusOK, Message: "Login", Token: tokenString}
}

func GetRefreshToken(userId string) (string, error) {
	user, err := searchUserById(userId)
	if err != nil {
		return "", err
	}
	if user.Name == "" {
		return "", fmt.Errorf("User with id: %s not found", userId)
	}
	return user.RefreshToken, nil

}

func Register(username string, password string) *structures.JSONMessage {
	user, err := searchUserByName(username)
	if err != nil {
		return &structures.JSONMessage{Status: http.StatusInternalServerError, Message: err.Error()}
	}
	if user.Id != "" {
		return &structures.JSONMessage{Status: http.StatusConflict, Message: "User already exists"}
	}
	id, err := uuid.NewRandom()

	if err != nil {
		log.Println(err.Error())
		return &structures.JSONMessage{Status: http.StatusInternalServerError, Message: err.Error()}
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	if err != nil {
		log.Println(err.Error())
		return &structures.JSONMessage{Status: http.StatusInternalServerError, Message: err.Error()}
	}

	refreshToken, err := tokenService.GenerateRefreshToken(id.String())
	if err != nil {
		return &structures.JSONMessage{Status: http.StatusInternalServerError, Message: err.Error()}
	}
	_, err = network.Pool.Exec(context.Background(), "INSERT INTO users (id, username, hash_password, refresh_token) VALUES ($1, $2, $3, $4)", id.String(), username, hashPassword, refreshToken)
	if err != nil {
		log.Println(err.Error())
		return &structures.JSONMessage{Status: http.StatusInternalServerError, Message: err.Error()}
	}

	tokenString, err := tokenService.GenerateAccessToken(id.String())
	if err != nil {
		return &structures.JSONMessage{Status: http.StatusInternalServerError, Message: err.Error()}
	}

	return &structures.JSONMessage{Status: http.StatusOK, Message: "Login", Token: tokenString}
}
