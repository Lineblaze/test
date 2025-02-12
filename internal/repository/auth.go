package repository

import (
	"avito_test/internal/entity"
	"avito_test/pkg/storage/postgres"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"time"
)

type AuthCache interface {
	Get(username string) (string, bool)
	Set(user string)
	Delete(username string)
}

type Auth struct {
	db postgres.Postgres
}

func NewAuth(db postgres.Postgres) Auth {
	return Auth{
		db: db,
	}
}

func (a Auth) Auth(ctx context.Context, auth entity.Auth) (*entity.Auth, error) {
	auth = entity.Auth{
		Id:       uuid.New(),
		Username: auth.Username,
		Password: auth.Password,
	}

	user := entity.User{
		Id:        auth.Id,
		Username:  auth.Username,
		Coin:      entity.Coin,
		CreatedAt: time.Now(),
	}

	err := postgres.ExecTx(ctx, a.db, func(tx postgres.Tx) error {
		queryAuth := `
			INSERT INTO auth (id, username, password) 
			VALUES ($1, $2, $3)
		`
		_, err := tx.Exec(ctx, queryAuth, auth.Id, auth.Username, auth.Password)
		if err != nil {
			return errors.Wrap(err, "failed to create auth in database")
		}

		queryUser := `
			INSERT INTO users (id, username, coin, created_at) 
			VALUES ($1, $2, $3, $4)
		`
		_, err = tx.Exec(ctx, queryUser, user.Id, user.Username, user.Coin, user.CreatedAt)
		if err != nil {
			return errors.Wrap(err, "failed to create user in database")
		}

		return nil
	})

	if err != nil {
		return nil, errors.Wrap(err, "transaction failed")
	}

	return &auth, nil
}
