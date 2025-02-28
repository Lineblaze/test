package service

import (
	"avito_test/internal/domain"
	"avito_test/internal/entity"
	"context"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"regexp"
)

type TransactionRepository interface {
	BuyItem(ctx context.Context, userID uuid.UUID, itemType string) error
	SendCoin(ctx context.Context, userID uuid.UUID, send entity.SendCoin) error
	GetInfo(ctx context.Context, userID uuid.UUID) (*entity.Info, error)
}

type Transaction struct {
	repo TransactionRepository
}

func NewTransaction(repo TransactionRepository) Transaction {
	return Transaction{
		repo: repo,
	}
}

func (t Transaction) Buy(ctx context.Context, userIDStr string, itemType string) error {
	if !validateUUID(userIDStr) {
		return domain.ErrInvalidCredentials
	}

	if !validateItemType(itemType) {
		return domain.ErrInvalidCredentials
	}

	userID, _ := uuid.Parse(userIDStr)

	err := t.repo.BuyItem(ctx, userID, itemType)
	if err != nil {
		return errors.Wrap(err, "failed to buy item")
	}

	return nil
}

func (t Transaction) Send(ctx context.Context, userIDStr string, req domain.SendCoinRequest) error {
	if !validateUUID(userIDStr) {
		return domain.ErrInvalidCredentials
	}

	userID, _ := uuid.Parse(userIDStr)

	entitySendCoin := entity.SendCoin{
		ToUser: req.ToUser,
		Amount: req.Amount,
	}

	err := t.repo.SendCoin(ctx, userID, entitySendCoin)
	if err != nil {
		return errors.Wrap(err, "failed to send coin")
	}

	return nil
}

func (t Transaction) Info(ctx context.Context, userIDStr string) (*domain.InfoResponse, error) {
	if !validateUUID(userIDStr) {
		return nil, domain.ErrInvalidCredentials
	}

	userID, _ := uuid.Parse(userIDStr)

	info, err := t.repo.GetInfo(ctx, userID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user info")
	}

	inventory := make([]domain.Item, 0, len(info.Inventory))
	receivedTransactions := make([]domain.CoinTransaction, 0, len(info.CoinHistory.Received))
	sentTransactions := make([]domain.CoinTransaction, 0, len(info.CoinHistory.Sent))

	for _, item := range info.Inventory {
		inventory = append(inventory, domain.Item{
			Type:     item.Type,
			Quantity: item.Quantity,
		})
	}

	for _, tx := range info.CoinHistory.Received {
		receivedTransactions = append(receivedTransactions, domain.CoinTransaction{
			FromUser: tx.FromUser,
			Amount:   tx.Amount,
		})
	}

	for _, tx := range info.CoinHistory.Sent {
		sentTransactions = append(sentTransactions, domain.CoinTransaction{
			ToUser: tx.ToUser,
			Amount: tx.Amount,
		})
	}

	res := domain.InfoResponse{
		Coins:     info.Coins,
		Inventory: inventory,
		CoinHistory: domain.CoinHistory{
			Received: receivedTransactions,
			Sent:     sentTransactions,
		},
	}

	return &res, nil
}

func validateUUID(userIDStr string) bool {
	_, err := uuid.Parse(userIDStr)
	return err == nil
}

func validateItemType(itemType string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9_-]{1,32}$`)
	return re.MatchString(itemType)
}
