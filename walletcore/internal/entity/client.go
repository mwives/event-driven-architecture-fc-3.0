package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidEntity         = errors.New("invalid entity")
	ErrClientAccountMismatch = errors.New("client and account mismatch")
)

type Client struct {
	ID        string
	Name      string
	Email     string
	Accounts  []*Account
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewClient(name, email string) (*Client, error) {
	client := &Client{
		ID:        uuid.New().String(),
		Name:      name,
		Email:     email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := client.Validate()
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (c *Client) Validate() error {
	if c.Name == "" || c.Email == "" {
		return ErrInvalidEntity
	}

	return nil
}

func (c *Client) Update(name, email string) error {
	tempClient := &Client{
		Name:  name,
		Email: email,
	}

	err := tempClient.Validate()
	if err != nil {
		return err
	}

	c.Name = name
	c.Email = email
	c.UpdatedAt = time.Now()

	return nil
}

func (c *Client) AddAccount(account *Account) error {
	if account.Client.ID != c.ID {
		return ErrClientAccountMismatch
	}
	c.Accounts = append(c.Accounts, account)
	return nil
}
