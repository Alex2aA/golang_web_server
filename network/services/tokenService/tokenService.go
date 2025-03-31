package tokenService

import (
	"github.com/golang-jwt/jwt/v5"
	"log"
	"os"
	"time"
)

func CreateToken(userID string, username string) (string, error) {
	claims := jwt.MapClaims{}
	claims["user_id"] = userID
	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Hour * 15).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))

	if err != nil {
		log.Println(err.Error())
		return "", err
	}

	return tokenString, nil
}
