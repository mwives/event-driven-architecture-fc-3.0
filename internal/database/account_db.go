package database

import (
	"database/sql"

	"github.com/mwives/microservices-fc-walletcore/internal/entity"
)

type AccountDB struct {
	DB *sql.DB
}

func NewAccountDB(db *sql.DB) *AccountDB {
	return &AccountDB{
		DB: db,
	}
}

func (a *AccountDB) FindByID(ID string) (*entity.Account, error) {
	var (
		account entity.Account
		client  entity.Client
	)
	account.Client = &client

	err := a.DB.QueryRow(`
		SELECT 
			a.id, a.client_id, a.balance, a.created_at, 
			c.id, c.name, c.email, c.created_at 
		FROM 
			accounts a 
		JOIN 
			clients c ON a.client_id = c.id 
		WHERE 
			a.id = ?`, ID).
		Scan(
			&account.ID, &account.Client.ID, &account.Balance, &account.CreatedAt,
			&account.Client.ID, &account.Client.Name, &account.Client.Email, &account.Client.CreatedAt,
		)

	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (a *AccountDB) Create(account *entity.Account) error {
	stmt, err := a.DB.Prepare("INSERT INTO accounts (id, client_id, balance, created_at) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(account.ID, account.Client.ID, account.Balance, account.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}
