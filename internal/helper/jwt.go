package helper

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	jwtSecret []byte
	jwtExpiry time.Duration
)

func SetJWTConfig(secret string, expiry time.Duration) {

	jwtSecret = []byte(secret)
	jwtExpiry = expiry
}

func GenerateToken(userID string) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"data": map[string]string{
			"user_id": userID,
		},
		"exp": time.Now().Add(jwtExpiry).Unix(),
		"iat": time.Now().Unix(),
	})

	return claims.SignedString(jwtSecret)
}

func ParseToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid or expired token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok && token.Valid {
		return nil, errors.New("invalid token claims")
	}
	return claims, nil
}
