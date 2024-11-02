package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrAmountMustBeGreaterThanZero = errors.New("amount must be greater than zero")
	ErrInsufficientBalance         = errors.New("insufficient balance")
)

type Transaction struct {
	ID          string
	AccountFrom *Account
	AccountTo   *Account
	Amount      float64
	CreatedAt   time.Time
}

func NewTransaction(accountFrom, accountTo *Account, amount float64) (*Transaction, error) {
	transaction := &Transaction{
		ID:          uuid.New().String(),
		AccountFrom: accountFrom,
		AccountTo:   accountTo,
		Amount:      amount,
		CreatedAt:   time.Now(),
	}

	err := transaction.Validate()
	if err != nil {
		return nil, err
	}

	transaction.Execute()

	return transaction, nil
}

func (t *Transaction) Validate() error {
	if t.Amount <= 0 {
		return ErrAmountMustBeGreaterThanZero
	}
	if t.AccountFrom.Balance < t.Amount {
		return ErrInsufficientBalance
	}
	return nil
}

func (t *Transaction) Execute() {
	t.AccountFrom.Debit(t.Amount)
	t.AccountTo.Credit(t.Amount)
}
