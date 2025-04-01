package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"idk_web_server001/network/httpHandlers"
	"idk_web_server001/structures"
	"log"
	"net/http"
	"os"
	"strings"
)

func parseToken(tokenString string, w http.ResponseWriter, r *http.Request) (*jwt.Token, error) {

	if tokenString == "" {
		httpHandlers.SendJSONMessage(w, structures.JSONMessage{Status: http.StatusUnauthorized, Message: "You are not authorized"})
		return nil, errors.New("You are not authorized")
	}

	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {

			httpHandlers.SendJSONMessage(w, structures.JSONMessage{Status: http.StatusUnauthorized, Message: "Unexpected signing method"})

			return nil, fmt.Errorf("Unexpected signing method")
		}
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	return token, err
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		token, err := parseToken(tokenString, w, r)
		if err != nil {
			log.Println("i'm there")
			if tokenString, err = httpHandlers.RefreshHandler(w, r); err != nil {
				return
			}
			token, err = parseToken(tokenString, w, r)

			if err != nil {
				httpHandlers.SendJSONMessage(w, structures.JSONMessage{Status: http.StatusInternalServerError, Message: err.Error()})
				log.Println(err.Error())
				return
			}
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
