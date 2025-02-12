package entity

import (
	"github.com/google/uuid"
	"time"
)

const Coin = 1000

type User struct {
	Id        uuid.UUID
	Username  string
	Coin      int64
	CreatedAt time.Time
}
