package internal

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// engine api check the range from iat,
// so it includes, iat & HS256
//   - required secret size is 32
func GenerateJws(secret []byte) (string, error) {
	if len(secret) != 32 {
		return "", errors.New("invalid secret size: expected 32 bytes")
	}

	iat := time.Now().Unix()
	claims := jwt.MapClaims{
		"iat": iat,
		"exp": iat + 60,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT: %w", err)
	}

	return tokenString, nil
}
