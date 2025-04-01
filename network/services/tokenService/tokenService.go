package tokenService

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"os"
	"time"
)

func GenerateTokens(userID string) (string, string, error) {
	claims := jwt.MapClaims{}
	claims["user_id"] = userID
	claims["exp"] = time.Now().Add(time.Second * 60).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))

	if err != nil {
		log.Println(err.Error())
		return "", "", err
	}

	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	refreshTokenString, err := refreshToken.SignedString([]byte(os.Getenv("SECRET_KEY")))

	if err != nil {
		return "", "", err
	}

	return refreshTokenString, tokenString, nil
}

func RefreshTokens(refreshTokenString string) (string, string, error) {
	refreshToken, err := jwt.Parse(refreshTokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {

			return nil, fmt.Errorf("Unexpected signing method")

		}
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil {
		log.Println(err.Error())
		return "", "", errors.New("Error in signing token")
	}

	if !refreshToken.Valid {
		return "", "", errors.New("Invalid refresh token")
	}

	claims, ok := refreshToken.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", errors.New("Token cast error")
	}

	userId, ok := claims["user_id"].(string)
	if !ok {
		return "", "", errors.New("Token type data error")
	}

	return GenerateTokens(userId)
}
