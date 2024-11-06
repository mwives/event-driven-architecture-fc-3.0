package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTransaction(t *testing.T) {
	client1, _ := NewClient("any_name1", "any_email1")
	account1 := NewAccount(client1)
	account1.Credit(1000)

	client2, _ := NewClient("any_name2", "any_email2")
	account2 := NewAccount(client2)
	account2.Credit(1000)

	transaction, err := NewTransaction(account1, account2, 100)
	assert.Nil(t, err)

	assert.Equal(t, transaction.AccountFrom.ID, account1.ID)
	assert.Equal(t, transaction.AccountTo.ID, account2.ID)
	assert.Equal(t, transaction.Amount, 100.0)
	assert.Equal(t, account1.Balance, 900.0)
	assert.Equal(t, account2.Balance, 1100.0)
}

func TestNewTransaction_WithInsufficientBalance(t *testing.T) {
	client1, _ := NewClient("any_name1", "any_email1")
	account1 := NewAccount(client1)
	account1.Credit(1000)

	client2, _ := NewClient("any_name2", "any_email2")
	account2 := NewAccount(client2)
	account2.Credit(1000)

	_, err := NewTransaction(account1, account2, 2000)
	assert.Equal(t, err, ErrInsufficientBalance)
}
