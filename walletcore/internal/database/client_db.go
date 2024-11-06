package database

import (
	"database/sql"
	"errors"
	"time"

	"github.com/mwives/event-driven-architecture-fc-3.0/walletcore/internal/entity"
)

type ClientDB struct {
	DB *sql.DB
}

func NewClientDB(db *sql.DB) *ClientDB {
	return &ClientDB{
		DB: db,
	}
}

func (c *ClientDB) FindByID(ID string) (*entity.Client, error) {
	client := &entity.Client{}
	var createdAtBytes []byte

	err := c.DB.
		QueryRow("SELECT id, name, email, created_at FROM clients WHERE id = ?", ID).
		Scan(&client.ID, &client.Name, &client.Email, &createdAtBytes)
	if err != nil {
		return nil, err
	}

	createdAt, err := time.Parse("2006-01-02 15:04:05", string(createdAtBytes))
	if err != nil {
		return nil, errors.New("error parsing created_at column to time.Time")
	}
	client.CreatedAt = createdAt

	return client, nil
}

func (c *ClientDB) Create(client *entity.Client) error {
	stmt, err := c.DB.Prepare("INSERT INTO clients (id, name, email, created_at) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(client.ID, client.Name, client.Email, client.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}
