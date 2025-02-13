package service

import (
	"avito_test/internal/domain"
	"avito_test/internal/entity"
	"avito_test/internal/jwt"
	"context"
	"github.com/pkg/errors"
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
