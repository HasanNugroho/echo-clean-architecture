package helper

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	redispkg "github.com/redis/go-redis/v9"
)

var (
	jwtSecret             []byte
	jwtExpiry             time.Duration
	jwtRefreshTokenExpiry time.Duration
	redisClient           *redispkg.Client
)

func SetJWTHelper(secret string, expiry time.Duration, refreshTokenExpiry time.Duration, redis *redispkg.Client) {
	jwtSecret = []byte(secret)
	jwtExpiry = expiry
	jwtRefreshTokenExpiry = refreshTokenExpiry
	redisClient = redis
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

func GenerateRefreshToken(userID string) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"type":    "refresh",
		"exp":     time.Now().Add(jwtRefreshTokenExpiry).Unix(),
		"iat":     time.Now().Unix(),
		"nbf":     time.Now().Add(jwtExpiry).Unix(),
	})

	return claims.SignedString(jwtSecret)
}

func ParseToken(tokenStr string) (jwt.MapClaims, error) {
	if IsTokenRevoked("refreshtoken:blacklist:"+tokenStr) || IsTokenRevoked("token:blacklist:"+tokenStr) {
		return nil, errors.New("invalid or expired token")
	}

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

func RevokeToken(tokenString string) error {
	ctx := context.Background()

	token, _, err := jwt.NewParser().ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return errors.New("failed to parse token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("invalid token claims")
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return errors.New("invalid expiration claim")
	}

	ttl := time.Until(time.Unix(int64(exp), 0))
	if err = redisClient.Set(ctx, "token:blacklist:"+tokenString, "revoked", ttl).Err(); err != nil {
		return errors.New("failed to store token in blacklist")
	}

	return nil
}

func RevokeRequestToken(refreshToken string) error {
	ctx := context.Background()

	token, _, err := jwt.NewParser().ParseUnverified(refreshToken, jwt.MapClaims{})
	if err != nil {
		return errors.New("failed to parse token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("invalid token claims")
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return errors.New("invalid expiration claim")
	}

	ttl := time.Until(time.Unix(int64(exp), 0))
	if err = redisClient.Set(ctx, "refreshtoken:blacklist:"+refreshToken, "revoked", ttl).Err(); err != nil {
		return errors.New("failed to store refresh token in blacklist")
	}

	return nil
}

func IsTokenRevoked(tokenString string) bool {
	_, err := redisClient.Get(context.Background(), tokenString).Result()
	return err == nil
}
