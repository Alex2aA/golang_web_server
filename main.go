package main

import (
	"github.com/joho/godotenv"
	"golang_web_server/network"
	"golang_web_server/network/router"
	"log"
	"net/http"
	"os"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	log.Println("Starting server...")

	network.DbConnect()

	r := router.SetupRouter()

	if err := os.MkdirAll(os.Getenv("VOLUME_USER_FILES"), os.ModePerm); err != nil {
		log.Fatal("Error creating directory", err.Error())
	}

	log.Fatal(http.ListenAndServe(":8080", r))
}
