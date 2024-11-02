package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAccount(t *testing.T) {
	client, _ := NewClient("any_name", "any_email")
	account := NewAccount(client)

	assert.NotNil(t, account)
	assert.Equal(t, client.ID, account.Client.ID)
}

func TestNewAccount_WithNilClient(t *testing.T) {
	account := NewAccount(nil)

	assert.Nil(t, account)
}

func TestAccount_Credit(t *testing.T) {
	client, _ := NewClient("any_name", "any_email")
	account := NewAccount(client)

	account.Credit(100)

	assert.Equal(t, 100.0, account.Balance)
}

func TestAccount_Debit(t *testing.T) {
	client, _ := NewClient("any_name", "any_email")
	account := NewAccount(client)

	account.Credit(100)
	account.Debit(50)

	assert.Equal(t, 50.0, account.Balance)
}
