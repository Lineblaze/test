package jwt

import (
	"avito_test/internal/config"
	"avito_test/pkg/logger"
	"github.com/pkg/errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const tokenExpiration = 24 * time.Hour

type Service struct {
	config *config.Config
	logger *logger.ApiLogger
}

type Claims struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

func NewJWTService(cfg *config.Config) *Service {
	return &Service{
		config: cfg,
	}
}

func (s *Service) GenerateJWT(claims Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       claims.ID,
		"username": claims.Username,
		"exp":      time.Now().Add(tokenExpiration).Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.config.Auth.Secret))
	if err != nil {
		return "", errors.WithMessage(err, "failed to sign JWT token")
	}

	return tokenString, nil
}

func (s *Service) ParseToken(tokenString string) (Claims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.config.Auth.Secret), nil
	})
	if err != nil {
		return Claims{}, errors.Wrap(err, "failed to parse token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return Claims{}, errors.New("invalid token claims")
	}

	userID, ok := claims["id"].(string)
	if !ok {
		return Claims{}, errors.New("invalid user ID in token")
	}

	username, _ := claims["username"].(string)

	return Claims{
		ID:       userID,
		Username: username,
	}, nil
}
