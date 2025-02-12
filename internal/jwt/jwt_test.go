package jwt

import (
	"avito_test/internal/config"
	"os"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

func NewMockConfig() *config.Config {
	return &config.Config{
		Auth: struct {
			Secret string `json:"secret"`
		}{
			Secret: os.Getenv("JWT_SECRET_KEY"),
		},
	}
}

func TestGenerateJWT(t *testing.T) {
	mockConfig := NewMockConfig()
	jwtService := NewJWTService(mockConfig)

	claims := Claims{
		ID:       "123",
		Username: "testuser",
	}

	token, err := jwtService.GenerateJWT(claims)
	assert.NoError(t, err, "GenerateJWT should not return an error")
	assert.NotEmpty(t, token, "Generated token should not be empty")

	parsedClaims, err := jwtService.ParseToken(token)
	assert.NoError(t, err, "ParseToken should not return an error")
	assert.Equal(t, claims.ID, parsedClaims.ID, "Parsed ID should match the original ID")
	assert.Equal(t, claims.Username, parsedClaims.Username, "Parsed username should match the original username")
}

func TestParseToken(t *testing.T) {
	mockConfig := NewMockConfig()
	jwtService := NewJWTService(mockConfig)

	claims := Claims{
		ID:       "123",
		Username: "testuser",
	}

	token, err := jwtService.GenerateJWT(claims)
	assert.NoError(t, err, "GenerateJWT should not return an error")

	parsedClaims, err := jwtService.ParseToken(token)
	assert.NoError(t, err, "ParseToken should not return an error")
	assert.Equal(t, claims.ID, parsedClaims.ID, "Parsed ID should match the original ID")
	assert.Equal(t, claims.Username, parsedClaims.Username, "Parsed username should match the original username")

	_, err = jwtService.ParseToken("invalid-token")
	assert.Error(t, err, "ParseToken should return an error for invalid token")

	expiredToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       claims.ID,
		"username": claims.Username,
		"exp":      time.Now().Add(-1 * time.Hour).Unix(),
	})
	expiredTokenString, err := expiredToken.SignedString([]byte(mockConfig.Auth.Secret))
	assert.NoError(t, err, "Failed to sign expired token")

	_, err = jwtService.ParseToken(expiredTokenString)
	assert.Error(t, err, "ParseToken should return an error for expired token")
}

func TestParseToken_InvalidClaims(t *testing.T) {
	mockConfig := NewMockConfig()
	jwtService := NewJWTService(mockConfig)

	invalidClaimsToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       123,
		"username": "testuser",
		"exp":      time.Now().Add(tokenExpiration).Unix(),
	})
	invalidClaimsTokenString, err := invalidClaimsToken.SignedString([]byte(mockConfig.Auth.Secret))
	assert.NoError(t, err, "Failed to sign token with invalid claims")

	_, err = jwtService.ParseToken(invalidClaimsTokenString)
	assert.Error(t, err, "ParseToken should return an error for invalid claims")
}
