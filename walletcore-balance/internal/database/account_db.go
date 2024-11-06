package database

import (
	"database/sql"

	"github.com/mwives/event-driven-architecture-fc-3.0/walletcore-balance/internal/entity"
)

type AccountDB struct {
	DB *sql.DB
}

func NewAccountDB(db *sql.DB) *AccountDB {
	return &AccountDB{
		DB: db,
	}
}

func (a *AccountDB) FindByID(id string) (*entity.Account, error) {
	var account entity.Account

	err := a.DB.
		QueryRow("SELECT id, balance, created_at, updated_at FROM accounts").
		Scan(&account.ID, &account.Balance, &account.CreatedAt, &account.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (a *AccountDB) Create(account *entity.Account) error {
	stmt, err := a.DB.Prepare("INSERT INTO accounts (id, balance, created_at, updated_at) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(account.ID, account.Balance, account.CreatedAt, account.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (a *AccountDB) Update(account *entity.Account) error {
	stmt, err := a.DB.Prepare("UPDATE accounts SET balance = ?, updated_at = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(account.Balance, account.UpdatedAt, account.ID)
	if err != nil {
		return err
	}

	return nil
}
