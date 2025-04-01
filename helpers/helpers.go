package helpers

import (
	"errors"
	"golang_web_server/structures"
	"log"
	"net/http"
	"os"
)

func CreateDirIfNotExist(r *http.Request) (*structures.JSONMessage, error) {
	userId, ok := r.Context().Value("userId").(string)

	var err error

	if !ok {
		return &structures.JSONMessage{Status: http.StatusInternalServerError, Message: "userId not found in context"}, errors.New("userId not found in context")
	}

	if err = os.MkdirAll(os.Getenv("VOLUME_USER_FILES")+"/"+userId, os.ModePerm); err != nil {
		log.Println("Error creating user directory", err.Error())
		return &structures.JSONMessage{Status: http.StatusInternalServerError, Message: "Error creating user directory"}, err
	}

	return &structures.JSONMessage{Status: http.StatusCreated, Message: "Created"}, nil
}

func FileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
