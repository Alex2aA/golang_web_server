package tokenService

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"os"
	"strings"
	"time"
)

func GenerateTokens(userID string, expired bool) (string, string, error) {
	claims := jwt.MapClaims{}
	claims["user_id"] = userID
	claims["exp"] = time.Now().Add(time.Second * 60).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))

	if err != nil {
		log.Println(err.Error())
		return "", "", err
	}

	if expired {
		claims["exp"] = time.Now().Add(time.Second * 120).Unix()

		refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		refreshTokenString, err := refreshToken.SignedString([]byte(os.Getenv("SECRET_KEY")))

		if err != nil {
			return "", "", err
		}

		return refreshTokenString, tokenString, nil
	}

	return "", tokenString, nil
}

func RefreshTokens(refreshTokenString string) (string, string, error) {
	refreshToken, err, expired := ParseToken(refreshTokenString)
	if err != nil && !expired {
		return "", "", err
	}

	if !refreshToken.Valid {
		return "", "", fmt.Errorf("refresh token is invalid")
	}

	claims, ok := refreshToken.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", errors.New("Token cast error")
	}

	userId, ok := claims["user_id"].(string)
	if !ok {
		return "", "", errors.New("Token type data error")
	}

	return GenerateTokens(userId, expired)
}

func IsExpired(token *jwt.Token) (bool, error) {

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false, fmt.Errorf("Cannot get claims from token")
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return false, fmt.Errorf("token exp field data error")
	}

	expirationTime := time.Unix(int64(exp), 0)
	return expirationTime.Before(time.Now()), nil
}

func ParseToken(tokenString string) (*jwt.Token, error, bool) {

	if tokenString == "" {
		return nil, errors.New("You are not authorized"), false
	}

	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method")
		}
		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	b, err := IsExpired(token)

	if err != nil {
		return nil, err, false
	}

	if b {
		return nil, fmt.Errorf("Token is expired"), true
	}
	return token, err, false
}
