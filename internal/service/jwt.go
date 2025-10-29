package service

import (
	"time"

	"github.com/Aleksey170999/go-loyaty/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

const (
	jwtTokenTTL = time.Hour * 24
)

type JWTService struct {
	secret []byte
}

type tokenClaims struct {
	UserID   int    `json:"user_id"`
	UserName string `json:"username"`
	jwt.RegisteredClaims
}

func NewJWTService(secret string) *JWTService {
	return &JWTService{
		secret: []byte(secret),
	}
}

func (s *JWTService) GenerateJWT(user models.User) (string, error) {
	claims := tokenClaims{
		UserID:   user.ID,
		UserName: user.UserName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(jwtTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

func (s *JWTService) ParseJWT(tokenStr string) (*tokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return s.secret, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}
