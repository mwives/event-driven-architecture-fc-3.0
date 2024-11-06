package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	name := "any_name"
	email := "any_email"

	client, err := NewClient(name, email)
	assert.Nil(t, err)

	assert.NotEmpty(t, client.ID)
	assert.Equal(t, name, client.Name)
	assert.Equal(t, email, client.Email)
	assert.NotZero(t, client.CreatedAt)
	assert.NotZero(t, client.UpdatedAt)
}

func TestNewClient_InvalidArgs(t *testing.T) {
	client, err := NewClient("", "")
	assert.NotNil(t, err)
	assert.Nil(t, client)
}

func TestUpdateClient(t *testing.T) {
	name := "any_name"
	email := "any_email"

	client, _ := NewClient(name, email)

	newName := "new_name"
	newEmail := "new_email"

	err := client.Update(newName, newEmail)
	assert.Nil(t, err)

	assert.Equal(t, newName, client.Name)
	assert.Equal(t, newEmail, client.Email)
	assert.NotZero(t, client.UpdatedAt)
}

func TestAddAccount(t *testing.T) {
	client, _ := NewClient("any_name", "any_email")
	account := NewAccount(client)

	err := client.AddAccount(account)
	assert.Nil(t, err)

	assert.Len(t, client.Accounts, 1)
	assert.Equal(t, account, client.Accounts[0])
}
