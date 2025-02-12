package entity

import (
	"github.com/google/uuid"
	"time"
)

type Auth struct {
	Id        uuid.UUID
	Username  string
	Password  string
	CreatedAt time.Time
}
