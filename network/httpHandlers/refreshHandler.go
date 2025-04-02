package httpHandlers

import (
	"errors"
	"golang_web_server/network/services/tokenService"
	"golang_web_server/structures"
	"log"
	"net/http"
	"strings"
)

func RefreshHandler(w http.ResponseWriter, r *http.Request) (string, error) {
	refreshTokenString := r.Header.Get("Refresh_authorization")
	if refreshTokenString == "" {
		return "", errors.New("you are not authorized")
	}

	refreshTokenString = strings.TrimPrefix(refreshTokenString, "Bearer ")
	newRefreshTokenString, tokenString, err := tokenService.RefreshTokens(refreshTokenString)
	if err != nil {
		return "", err
	}

	token, err, expired := tokenService.ParseToken(refreshTokenString)

	if token == nil {
		log.Println("tokenService.ParseToken err")
	}

	if err != nil {
		return "", err
	}

	if expired {
		SendJSONMessage(w, structures.JSONMessage{Status: http.StatusOK, Message: "Token refresh", Token: tokenString, RefreshToken: newRefreshTokenString})
		return tokenString, nil
	}

	SendJSONMessage(w, structures.JSONMessage{Status: http.StatusOK, Message: "Token refresh", Token: tokenString})

	return tokenString, nil
}
