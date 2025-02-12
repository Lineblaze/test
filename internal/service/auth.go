package service

import (
	"avito_test/internal/domain"
	"avito_test/internal/entity"
	"avito_test/internal/jwt"
	"context"
	"github.com/pkg/errors"
	"regexp"
)

type AuthRepository interface {
	Auth(ctx context.Context, user entity.Auth) (*entity.Auth, error)
}

type Auth struct {
	repo AuthRepository
	jwt  *jwt.Service
}

func NewAuth(repo AuthRepository, jwtService *jwt.Service) Auth {
	return Auth{
		repo: repo,
		jwt:  jwtService,
	}
}

func (a Auth) Auth(ctx context.Context, req domain.AuthRequest) (*domain.AuthResponse, error) {
	if !validateUsername(req.Username) || !validatePassword(req.Password) {
		return nil, domain.ErrInvalidCredentials
	}

	entityAuth := entity.Auth{
		Username: req.Username,
		Password: req.Password,
	}

	authUser, err := a.repo.Auth(ctx, entityAuth)
	if err != nil {
		return nil, errors.Wrap(err, "create user failed")
	}

	token, err := a.jwt.GenerateJWT(jwt.Claims{
		ID:       authUser.Id.String(),
		Username: authUser.Username,
	})
	if err != nil {
		return nil, domain.ErrUnauthorized
	}

	res := domain.AuthResponse{
		Token: token,
	}

	return &res, nil
}

func validateUsername(username string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._-]{3,32}$`)
	return re.MatchString(username)
}

func validatePassword(password string) bool {
	re := regexp.MustCompile(`^[A-Za-z\d!@#$%^&*()\-_+=]{8,}$`)
	hasLetter := regexp.MustCompile(`[A-Za-z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`\d`).MatchString(password)
	return re.MatchString(password) && hasLetter && hasDigit
}
