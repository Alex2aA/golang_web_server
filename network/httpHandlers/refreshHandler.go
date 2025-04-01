package httpHandlers

import (
	"errors"
	"golang_web_server/network/services/tokenService"
	"golang_web_server/structures"
	"net/http"
	"strings"
)

func RefreshHandler(w http.ResponseWriter, r *http.Request) (string, error) {
	refreshTokenString := r.Header.Get("Refresh_authorization")
	if refreshTokenString == "" {
		SendJSONMessage(w, structures.JSONMessage{Status: http.StatusUnauthorized, Message: "You are not authorized"})
		return "", errors.New("you are not authorized")
	}

	refreshTokenString = strings.TrimPrefix(refreshTokenString, "Bearer ")
	newRefreshTokenString, tokenString, err := tokenService.RefreshTokens(refreshTokenString)

	if err != nil {
		SendJSONMessage(w, structures.JSONMessage{Status: http.StatusInternalServerError, Message: err.Error()})
		return "", err
	}

	SendJSONMessage(w, structures.JSONMessage{Status: http.StatusOK, Message: "Token refresh", Token: tokenString, RefreshToken: newRefreshTokenString})
	return tokenString, nil
}
