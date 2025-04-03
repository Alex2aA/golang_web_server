package tokenService

import (
	"github.com/golang-jwt/jwt/v5"
	"golang_web_server/structures"
	"os"
	"time"
)

const (
	accessTokenExp  = time.Second * 30
	refreshTokenExp = time.Minute * 2
)

// ------------- Generate Tokens -------------
func GenerateAccessToken(userId string) (string, error) {
	claims := structures.TokenClaims{
		UserId:    userId,
		BlackList: false,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTokenExp)),
			Issuer:    "tokenService",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("SECRET_KEY")))
}

func GenerateRefreshToken(userId string) (string, error) {
	claims := structures.TokenClaims{
		UserId:    userId,
		BlackList: false,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshTokenExp)),
			Issuer:    "tokenService",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("SECRET_KEY")))
}

//------------- Generate Tokens -------------
