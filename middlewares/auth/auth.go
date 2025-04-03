package auth

import (
	"context"
	"golang_web_server/network/httpHandlers"
	"golang_web_server/network/services/tokenService/tokenParser"
	"golang_web_server/structures"
	"net/http"
	"strings"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tokenString := r.Header.Get("Authorization")
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		claims, err, msg := tokenParser.ParseToken(tokenString)

		if msg != nil {
			httpHandlers.SendJSONMessage(w, *msg)
		}

		if err != nil {
			httpHandlers.SendJSONMessage(w, structures.JSONMessage{Status: http.StatusUnauthorized, Message: err.Error()})
			return
		}

		ctx := context.WithValue(r.Context(), "userId", claims.UserId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
