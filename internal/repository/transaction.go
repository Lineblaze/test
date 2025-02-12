package repository

import (
	"avito_test/internal/entity"
	"avito_test/pkg/storage/postgres"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

type Transaction struct {
	db postgres.Postgres
}

func NewTransaction(db postgres.Postgres) Transaction {
	return Transaction{
		db: db,
	}
}

func (t Transaction) GetInfo(ctx context.Context, userID uuid.UUID) (*entity.Info, error) {
	var info entity.Info

	err := postgres.ExecTx(ctx, t.db, func(tx postgres.Tx) error {
		query := `SELECT coin FROM users WHERE id = $1`
		err := tx.Get(ctx, &info.Coins, query, userID)
		if err != nil {
			return errors.WithMessage(err, "failed to get user coins")
		}

		query = `SELECT type, quantity FROM user_items WHERE user_id = $1`
		err = tx.Select(ctx, &info.Inventory, query, userID)
		if err != nil {
			return errors.WithMessage(err, "failed to get user inventory")
		}

		query = `SELECT from_user, amount FROM coin_transactions WHERE to_user = (SELECT username FROM users WHERE id = $1)`
		err = tx.Select(ctx, &info.CoinHistory.Received, query, userID)
		if err != nil {
			return errors.WithMessage(err, "failed to get received transactions")
		}

		query = `SELECT to_user, amount FROM coin_transactions WHERE from_user = (SELECT username FROM users WHERE id = $1)`
		err = tx.Select(ctx, &info.CoinHistory.Sent, query, userID)
		if err != nil {
			return errors.WithMessage(err, "failed to get sent transactions")
		}

		return nil
	})

	if err != nil {
		return nil, errors.Wrap(err, "transaction failed")
	}

	return &info, nil
}

func (t Transaction) BuyItem(ctx context.Context, userID uuid.UUID, itemType string) error {
	err := postgres.ExecTx(ctx, t.db, func(tx postgres.Tx) error {
		var price int64
		query := `SELECT price 
				  FROM items 
				  WHERE type = $1 
				  LIMIT 1`
		err := tx.Get(ctx, &price, query, itemType)
		if err != nil {
			return errors.WithMessage(err, "failed to get item price")
		}

		var user entity.User
		query = `SELECT id, username, coin 
				 FROM users 
				 WHERE id = $1`
		err = tx.Get(ctx, &user, query, userID)
		if err != nil {
			return errors.WithMessage(err, "failed to get user by ID")
		}

		newCoinBalance := user.Coin - price
		query = `UPDATE users 
				 SET coin = $1 
				 WHERE id = $2 
				 RETURNING coin`
		err = tx.Get(ctx, &newCoinBalance, query, newCoinBalance, userID)
		if err != nil {
			return errors.WithMessage(err, "failed to update user coins")
		}

		query = `INSERT INTO user_items (user_id, type, quantity) 
				  VALUES ($1, $2, 1) 
				  ON CONFLICT (user_id, type) 
				  DO UPDATE SET quantity = user_items.quantity + 1`
		_, err = tx.Exec(ctx, query, userID, itemType)
		if err != nil {
			return errors.WithMessage(err, "failed to add user item")
		}

		return nil
	})

	if err != nil {
		return errors.Wrap(err, "transaction failed")
	}

	return nil
}

func (t Transaction) SendCoin(ctx context.Context, userID uuid.UUID, send entity.SendCoin) error {
	err := postgres.ExecTx(ctx, t.db, func(tx postgres.Tx) error {
		var sender struct {
			Username string
			Coin     int
		}

		query := `SELECT username, coin FROM users WHERE id = $1 FOR UPDATE`
		err := tx.Get(ctx, &sender, query, userID)
		if err != nil {
			return errors.WithMessage(err, "failed to get sender")
		}

		if sender.Coin < send.Amount {
			return errors.New("insufficient balance")
		}

		newSenderBalance := sender.Coin - send.Amount
		query = `UPDATE users SET coin = $1 WHERE id = $2 RETURNING coin`
		err = tx.Get(ctx, &newSenderBalance, query, newSenderBalance, userID)
		if err != nil {
			return errors.WithMessage(err, "failed to update sender balance")
		}

		query = `UPDATE users SET coin = coin + $1 WHERE username = $2 RETURNING coin`
		var newReceiverBalance int
		err = tx.Get(ctx, &newReceiverBalance, query, send.Amount, send.ToUser)
		if err != nil {
			return errors.WithMessage(err, "failed to update receiver balance")
		}

		query = `INSERT INTO coin_transactions (from_user, to_user, amount) VALUES ($1, $2, $3)`
		_, err = tx.Exec(ctx, query, sender.Username, send.ToUser, send.Amount)
		if err != nil {
			return errors.WithMessage(err, "failed to insert coin transaction")
		}

		return nil
	})

	if err != nil {
		return errors.Wrap(err, "transaction failed")
	}
	return nil
}
