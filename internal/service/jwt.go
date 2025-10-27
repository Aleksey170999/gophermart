package service

import (
	"time"

	"github.com/Aleksey170999/go-loyaty/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

const (
	jwtSecret   = "your_secret_key" // Замените на свой ключ
	jwtTokenTTL = time.Hour * 24
)

type tokenClaims struct {
	UserID   int    `json:"user_id"`
	UserName string `json:"username"`
	jwt.RegisteredClaims
}

func GenerateJWT(user models.User) (string, error) {
	claims := tokenClaims{
		UserID:   user.ID,
		UserName: user.UserName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(jwtTokenTTL)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

func ParseJWT(tokenStr string) (*tokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*tokenClaims)
	if !ok || !token.Valid {
		return nil, err
	}
	return claims, nil
}
