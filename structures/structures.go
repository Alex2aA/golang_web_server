package structures

import "github.com/golang-jwt/jwt/v5"

type JSONMessage struct {
	Status  int    `json:"status"`
	Message string `json:"message,omitempty"`
	Token   string `json:"token,omitempty"`
}

type User struct {
	Id           string `json:"id,omitempty"`
	Name         string `json:"username"`
	Password     string `json:"password"`
	RefreshToken string `json:"refreshToken,omitempty"`
}

type Post struct {
	Id          string `json:"id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Text        string `json:"text"`
}

type OptsFile struct {
	IsAvatar bool `json:"IsAvatar,omitempty"`
	IsImage  bool `json:"IsImage,omitempty"`
}

type TokenClaims struct {
	UserId    string `json:"userId"`
	BlackList bool   `json:"blackList"`
	jwt.RegisteredClaims
}
