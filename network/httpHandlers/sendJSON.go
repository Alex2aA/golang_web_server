package httpHandlers

import (
	"encoding/json"
	"idk_web_server001/structures"
	"log"
	"net/http"
)

func SendJSONMessage(w http.ResponseWriter, message structures.JSONMessage) {
	w.WriteHeader(message.Status)
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(message); err != nil {
		log.Println("Unable to encode message:", err.Error())
	}
}
