package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/thecipherdev/goauth/config"
)

type CustomClaims struct {
	Username  string `json:"username"`
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

func GenerateToken(sub, username, tokenType string, duration time.Duration) (string, error) {
	var (
		key   []byte
		token *jwt.Token
	)

	cfg := config.Get()
	key = []byte(cfg.JWTSecret)

	claims := jwt.RegisteredClaims{
		Issuer:    "my-auth-server",
		Subject:   sub,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
	}

	custom := &CustomClaims{
		TokenType:        tokenType,
		Username:         username,
		RegisteredClaims: claims,
	}

	token = jwt.NewWithClaims(jwt.SigningMethodHS256, custom)

	signedKey, err := token.SignedString(key)

	if err != nil {
		return "", fmt.Errorf("Failed to sign token: %w", err)
	}

	return signedKey, nil

}

func GenerateAccessToken() (string, error) {
	return GenerateToken("1231412312", "johndoe", "access", 15*time.Minute)
}

func GenerateRefreshToken() (string, error) {
	return GenerateToken("1231412312", "johndoe", "refresh", 7*24*time.Hour)
}

func ValidateToken(tokenStr string) (*CustomClaims, error) {
	cfg := config.Get()
	key := []byte(cfg.JWTSecret)

	claims := &CustomClaims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}
		return key, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("token expired: %w", err)
		}
		if errors.Is(err, jwt.ErrTokenNotValidYet) {
			return nil, fmt.Errorf("token not active yet: %w", err)
		}
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	if !token.Valid {
		return nil, errors.New("token is invalid")
	}

	return claims, nil
}
