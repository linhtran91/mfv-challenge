package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenClaims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

type jwtHMACBuilder struct {
	secretKey []byte
	duration  time.Duration
}

func NewjwtHMACBuilder(secret string, duration time.Duration) *jwtHMACBuilder {
	return &jwtHMACBuilder{
		secretKey: []byte(secret),
		duration:  duration,
	}
}

func (b *jwtHMACBuilder) Encode(userID int64) (string, error) {
	now := time.Now()
	claims := TokenClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(b.duration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(b.secretKey)
}

func (b *jwtHMACBuilder) Decode(tokenString string) (int64, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return b.secretKey, nil
	})
	if err != nil {
		return -1, err
	}
	if claims, ok := token.Claims.(*TokenClaims); ok {
		return claims.UserID, nil
	}
	return -1, errors.New("unknown claims type, cannot proceed")
}
