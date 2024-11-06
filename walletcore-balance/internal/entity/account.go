package entity

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID        string
	Balance   float64
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewAccount() *Account {
	return &Account{
		ID:        uuid.New().String(),
		Balance:   0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
