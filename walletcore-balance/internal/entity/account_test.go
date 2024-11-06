package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAccount(t *testing.T) {
	account := NewAccount()
	assert.NotNil(t, account)
	assert.NotEmpty(t, account.ID)
}
