package internal_test

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/poteto-go/go-alchemy-sdk/constant"
	"github.com/poteto-go/go-alchemy-sdk/internal"
	"github.com/stretchr/testify/assert"
)

func TestGenerateJws(t *testing.T) {
	t.Run("generate signed token (iat <= time w/ HS256)", func(t *testing.T) {
		// Arrange
		sampleHex := "bcdf20249abf0ed6d944c0288fad489e33f66b3960d9e6229c1cd214ed3bbe31"
		secret, _ := internal.DecodeHex(sampleHex)
		before := time.Now().Unix()

		// Act
		tokenStr, err := internal.GenerateJws(secret)
		claims := &jwt.MapClaims{}
		decoded, _ := jwt.ParseWithClaims(
			tokenStr,
			claims,
			func(t *jwt.Token) (any, error) {
				return "a", nil
			},
		)
		after := time.Now().Unix()

		// Assert
		assert.Nil(t, err)
		assert.Equal(t, decoded.Method, jwt.SigningMethodHS256)
		iat, _ := decoded.Claims.GetIssuedAt()
		assert.True(t, before <= iat.Unix())
		assert.True(t, after >= iat.Unix())
		exp, _ := decoded.Claims.GetExpirationTime()
		assert.Equal(t, iat.Unix()+constant.GethJwsIatWindowSec, exp.Unix())
	})

	t.Run("fail on invalid jwt secret", func(t *testing.T) {
		// Act
		_, err := internal.GenerateJws([]byte("invalid"))

		// Assert
		assert.Error(t, err)
	})
}
