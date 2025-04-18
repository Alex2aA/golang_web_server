package uploadAvatarService

import (
	"errors"
	"fmt"
	"golang_web_server/structures"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Максимальный размер загружаемого файла(аватарки)
const maxUploadSizeAvatar = 2 * 1024 * 1024

var (
	file           multipart.File
	header         *multipart.FileHeader
	newFileName    string
	avatarFilePath string
	userId         string
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

func checkAvatarIfExists(r *http.Request) (bool, error) {
	var ok bool
	userId, ok = r.Context().Value("userId").(string)
	if !ok {
		return false, errors.New("userId not found in context")
	}

	ext := filepath.Ext(header.Filename)

	newFileName = fmt.Sprintf("%s%s", userId, ext)

	avatarFilePath = filepath.Join(os.Getenv("VOLUME_USER_FILES")+"/"+userId+"/avatar", newFileName)

	return FileExists(avatarFilePath), nil
}

func InitUploadAvatarService(w http.ResponseWriter, r *http.Request) *structures.JSONMessage {
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSizeAvatar)

	var err error

	defer r.Body.Close()

	if err = r.ParseMultipartForm(maxUploadSizeAvatar); err != nil {
		return &structures.JSONMessage{Status: http.StatusBadRequest, Message: "File too big. Max size " + fmt.Sprintf("%d MB", maxUploadSizeAvatar/(1024*1024))}
	}
	file, header, err = r.FormFile("avatar")
	if err != nil {
		log.Println(err.Error())
		return &structures.JSONMessage{Status: http.StatusBadRequest, Message: "Invalid file"}
	}
	defer file.Close()

	var message *structures.JSONMessage

	if message, err = readFile(); err != nil {
		return message
	}
	if message, err = saveFile(r); err != nil {
		return message
	}
	return message
}

func readFile() (*structures.JSONMessage, error) {
	buff := make([]byte, 512)
	if _, err := file.Read(buff); err != nil {
		return &structures.JSONMessage{Status: http.StatusBadRequest, Message: "Invalid file content"}, err
	}
	if http.DetectContentType(buff) != "image/png" && http.DetectContentType(buff) != "image/jpeg" {
		return &structures.JSONMessage{Status: http.StatusBadRequest, Message: "Invalid file format"}, errors.New("Invalid file format")
	}

	return &structures.JSONMessage{Status: http.StatusOK}, nil
}

func getExtFile() (string, error) {
	files, err := os.ReadDir("./uploads/" + userId + "/avatar")

	if err != nil {
		log.Println(err.Error())
		return "", err
	}

	var ext string

	for _, file := range files {
		fileName := file.Name()
		ext = strings.TrimPrefix(fileName, filepath.Ext(fileName))
	}

	return ext, nil

}

func saveFile(r *http.Request) (*structures.JSONMessage, error) {

	var b bool
	var e error

	if b, e = checkAvatarIfExists(r); e != nil {
		return &structures.JSONMessage{Status: http.StatusInternalServerError, Message: "userId not found in context"}, e
	}

	if b {

		ext, err := getExtFile()
		if err != nil {
			log.Println(err.Error())
			return &structures.JSONMessage{Status: http.StatusInternalServerError, Message: "Invalid ext operation"}, err
		}

		filePathAvatar := filepath.Join(os.Getenv("VOLUME_USER_FILES") + "/" + userId + "/avatar" + "/" + ext)

		err = os.Remove(filePathAvatar)

		if err != nil {
			log.Println(err.Error())
			return &structures.JSONMessage{Status: http.StatusInternalServerError, Message: "Error while removing the file"}, e
		}
		log.Println("File removed")
	}

	message, err := CreateDirIfNotExist(r)
	if err != nil {
		return message, err
	}

	if err = os.MkdirAll(os.Getenv("VOLUME_USER_FILES")+"/"+userId+"/avatar", os.ModePerm); err != nil {
		log.Println("Error creating avatar user directory", err.Error())
		return &structures.JSONMessage{Status: http.StatusInternalServerError, Message: "Error creating avatar user directory"}, err
	}

	out, err := os.Create(avatarFilePath)
	if err != nil {
		log.Println("Error save file", err.Error())
		return &structures.JSONMessage{Status: http.StatusInternalServerError, Message: "Failed to save file"}, err
	}

	defer out.Close()

	if _, err = file.Seek(0, io.SeekStart); err != nil {
		log.Println("Error read file", err.Error())
		return &structures.JSONMessage{Status: http.StatusInternalServerError, Message: "Failed to read file"}, err
	}

	if _, err = io.Copy(out, file); err != nil {
		log.Println("Error save file", err.Error())
		return &structures.JSONMessage{Status: http.StatusInternalServerError, Message: "Failed to save file"}, err
	}

	return &structures.JSONMessage{Status: http.StatusCreated, Message: "File loaded"}, nil
}
