package tokenService

import (
	"errors"
	"fmt"
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

func ParseToken(tokenString string) (*structures.TokenClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&structures.TokenClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SECRET_KEY")), nil
		},
	)

	if err != nil && !errors.Is(err, jwt.ErrTokenExpired) {
		return nil, fmt.Errorf("parsing failed: %w", err)
	}

	claims, ok := token.Claims.(*structures.TokenClaims)
	if !ok {
		return &structures.TokenClaims{}, errors.New("invalid token claims structure")
	}

	if claims.BlackList {
		return &structures.TokenClaims{}, errors.New("invalid token, token in blacklist")
	}

	if isExpired(token, err) {
		claims.BlackList = true

		return claims, nil
	}

	if err != nil || !token.Valid {
		return &structures.TokenClaims{}, fmt.Errorf("invalid token: %w", err)
	}

	return claims, nil
}

func isExpired(token *jwt.Token, err error) bool {
	if !token.Valid && errors.Is(err, jwt.ErrTokenExpired) {
		return true
	}
	return false
}
