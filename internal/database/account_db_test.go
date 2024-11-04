package database

import (
	"database/sql"
	"testing"

	"github.com/mwives/microservices-fc-walletcore/internal/entity"
	"github.com/stretchr/testify/suite"
)

type AccountDbTestSuite struct {
	suite.Suite
	db        *sql.DB
	accountDB *AccountDB
}

func (s *AccountDbTestSuite) SetupSuite() {
	db, err := sql.Open("sqlite3", ":memory:")
	s.Nil(err)
	s.db = db

	db.Exec("CREATE TABLE accounts (id TEXT PRIMARY KEY, client_id TEXT, balance INTEGER, created_at DATETIME)")
	db.Exec("CREATE TABLE clients (id TEXT PRIMARY KEY, name TEXT, email TEXT, created_at DATETIME)")

	s.accountDB = NewAccountDB(db)
}

func (s *AccountDbTestSuite) TearDownSuite() {
	defer s.db.Close()
	s.db.Exec("DROP TABLE accounts")
	s.db.Exec("DROP TABLE clients")
}

func TestAccountDbTestSuite(t *testing.T) {
	suite.Run(t, new(AccountDbTestSuite))
}

func (s *AccountDbTestSuite) TestCreate() {
	client, _ := entity.NewClient("any_name", "any_email")
	account := entity.NewAccount(client)
	err := s.accountDB.Create(account)
	s.Nil(err)
}

func (s *AccountDbTestSuite) TestFindByID() {
	client, _ := entity.NewClient("any_name", "any_email")
	s.db.Exec(
		"INSERT INTO clients (id, name, email, created_at) VALUES (?, ?, ?, ?)",
		client.ID, client.Name, client.Email, client.CreatedAt,
	)

	account := entity.NewAccount(client)
	err := s.accountDB.Create(account)
	s.Nil(err)

	accountDB, err := s.accountDB.FindByID(account.ID)
	s.Nil(err)
	s.Equal(account.ID, accountDB.ID)
	s.Equal(account.Client.ID, accountDB.Client.ID)
	s.Equal(account.Balance, accountDB.Balance)
	s.Equal(account.Client.Name, accountDB.Client.Name)
	s.Equal(account.Client.Email, accountDB.Client.Email)
}
