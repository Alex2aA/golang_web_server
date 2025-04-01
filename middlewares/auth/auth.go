package auth

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"golang_web_server/network/httpHandlers"
	"golang_web_server/network/services/tokenService"
	"golang_web_server/structures"
	"log"
	"net/http"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		token, err, expired := tokenService.ParseToken(tokenString)
		if err != nil && !expired {
			httpHandlers.SendJSONMessage(w, structures.JSONMessage{Status: http.StatusUnauthorized, Message: err.Error()})
			return
		}

		if expired {
			tokenString, err = httpHandlers.RefreshHandler(w, r)
			if err != nil {
				httpHandlers.SendJSONMessage(w, structures.JSONMessage{Status: http.StatusUnauthorized, Message: err.Error()})
				return
			}
		}

		token, err, expired = tokenService.ParseToken(tokenString)
		if err != nil && !expired {
			httpHandlers.SendJSONMessage(w, structures.JSONMessage{Status: http.StatusUnauthorized, Message: err.Error()})
			return
		}

		if err != nil {
			httpHandlers.SendJSONMessage(w, structures.JSONMessage{Status: http.StatusInternalServerError, Message: err.Error()})
			log.Println(err.Error())
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			httpHandlers.SendJSONMessage(w, structures.JSONMessage{Status: http.StatusUnauthorized, Message: "Token cast error"})
			return
		}

		userId, ok := claims["user_id"].(string)
		if !ok {
			httpHandlers.SendJSONMessage(w, structures.JSONMessage{Status: http.StatusUnauthorized, Message: "Token type data error"})
			return
		}

		//username, ok := claims["username"].(string)
		//if !ok {
		//	httpHandlers.SendJSONMessage(w, structures.JSONMessage{Status: http.StatusUnauthorized, Message: "Token type data error"})
		//	return
		//}

		ctx := context.WithValue(r.Context(), "userId", userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
