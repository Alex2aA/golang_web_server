package auth

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"golang_web_server/network/httpHandlers"
	"golang_web_server/structures"
	"net/http"
	"os"
	"strings"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			w.WriteHeader(http.StatusUnauthorized)
			httpHandlers.SendJSONMessage(w, structures.JSONMessage{Status: http.StatusUnauthorized, Message: "You are not authorized"})
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {

				httpHandlers.SendJSONMessage(w, structures.JSONMessage{Status: http.StatusUnauthorized, Message: "Unexpected signing method"})

				return nil, fmt.Errorf("Unexpected signing method")
			}
			return []byte(os.Getenv("SECRET_KEY")), nil
		})
		if err != nil {
			httpHandlers.SendJSONMessage(w, structures.JSONMessage{Status: http.StatusUnauthorized, Message: "Error in signing token"})
			return
		}

		if !token.Valid {
			httpHandlers.SendJSONMessage(w, structures.JSONMessage{Status: http.StatusUnauthorized, Message: "Your token is expired"})
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
