package database

import (
	"database/sql"

	"github.com/mwives/microservices-fc-walletcore/internal/entity"
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

	err := c.DB.
		QueryRow("SELECT id, name, email, created_at FROM clients WHERE id = ?", ID).
		Scan(&client.ID, &client.Name, &client.Email, &client.CreatedAt)
	if err != nil {
		return nil, err
	}

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
