package httpHandlers

import (
	"context"
	"idk_web_server001/network"
	"idk_web_server001/structures"
	"log"
	"net/http"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	err := network.Pool.Ping(context.Background())
	if err != nil {
		log.Println("Unable to connect to pool", err.Error())
		SendJSONMessage(w, structures.JSONMessage{Status: http.StatusServiceUnavailable, Message: "Pool is down"})
		return
	}
	SendJSONMessage(w, structures.JSONMessage{Status: http.StatusOK, Message: "All ok"})
}
