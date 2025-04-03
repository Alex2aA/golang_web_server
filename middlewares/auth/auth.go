package auth

import (
	"context"
	"golang_web_server/network/httpHandlers"
	"golang_web_server/network/services/tokenService"
	"golang_web_server/structures"
	"net/http"
	"strings"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tokenString := r.Header.Get("Authorization")
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		claims, err := tokenService.ParseToken(tokenString)

		if err != nil {
			httpHandlers.SendJSONMessage(w, structures.JSONMessage{Status: http.StatusUnauthorized, Message: err.Error()})
			return
		}

		httpHandlers.SendJSONMessage(w, structures.JSONMessage{Status: http.StatusOK, Message: "Token ok", Token: tokenString})

		ctx := context.WithValue(r.Context(), "userId", claims.UserId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
