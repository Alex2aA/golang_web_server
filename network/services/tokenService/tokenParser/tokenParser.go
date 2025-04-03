package tokenParser

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"golang_web_server/network/services/tokenService"
	"golang_web_server/network/services/userService"
	"golang_web_server/structures"
	"net/http"
	"os"
	"time"
)

func ParseToken(tokenString string) (*structures.TokenClaims, error, *structures.JSONMessage) {

	token, err := jwt.ParseWithClaims(
		tokenString,
		&structures.TokenClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SECRET_KEY")), nil
		},
	)

	if err != nil && !errors.Is(err, jwt.ErrTokenExpired) {
		return nil, fmt.Errorf("parsing failed: %w", err), nil
	}

	claims, ok := token.Claims.(*structures.TokenClaims)
	if !ok {
		return &structures.TokenClaims{}, errors.New("invalid token claims structure"), nil
	}

	b, err := userService.CheckTokenBlackList(tokenString)
	if err != nil {
		return &structures.TokenClaims{}, err, nil
	}

	if b {
		return &structures.TokenClaims{}, errors.New("invalid token, token in blacklist"), nil
	}

	if isExpired(token, err) {
		newTokenString, err := generateTokenWRT(claims)
		if err != nil {
			return &structures.TokenClaims{}, err, nil
		}
		err = userService.AddTokenBlackList(tokenString)
		if err != nil {
			return &structures.TokenClaims{}, err, nil
		}
		newToken, err, _ := ParseToken(newTokenString)
		if err != nil {
			return &structures.TokenClaims{}, err, nil
		}

		return newToken, nil, &structures.JSONMessage{Status: http.StatusCreated, Message: "Token updated", Token: newTokenString}
	}

	return claims, nil, nil
}

func generateTokenWRT(claims *structures.TokenClaims) (string, error) {
	refreshTokenString, err := userService.GetRefreshToken(claims.UserId)
	if err != nil {
		return "", err
	}

	refreshToken, err := jwt.ParseWithClaims(
		refreshTokenString,
		&structures.TokenClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SECRET_KEY")), nil
		},
	)

	if isExpired(refreshToken, err) {
		userService.AddTokenBlackList(refreshTokenString)
		return "", errors.New("refresh token is expired")
	}

	newTokenString, err := tokenService.GenerateAccessToken(claims.UserId)
	if err != nil {
		return "", err
	}

	return newTokenString, nil
}

func isExpired(token *jwt.Token, err error) bool {
	exp, err := token.Claims.GetExpirationTime()
	if err != nil {
		return false
	}
	if exp.Before(time.Now()) {
		return true
	}

	return false
}
