package database

import (
	"database/sql"
	"testing"

	"github.com/mwives/microservices-fc-walletcore/internal/entity"
	"github.com/stretchr/testify/suite"
)

type TransactionDBTestSuite struct {
	suite.Suite
	db            *sql.DB
	TransactionDB *TransactionDB
	clientFrom    *entity.Client
	clientTo      *entity.Client
	accountFrom   *entity.Account
	accountTo     *entity.Account
}

func (s *TransactionDBTestSuite) SetupSuite() {
	db, err := sql.Open("sqlite3", ":memory:")
	s.Nil(err)
	s.db = db

	db.Exec("CREATE TABLE clients (id TEXT PRIMARY KEY, name TEXT, email TEXT, created_at DATETIME)")
	db.Exec("CREATE TABLE accounts (id TEXT PRIMARY KEY, client_id TEXT, balance INTEGER, created_at DATETIME)")
	db.Exec("CREATE TABLE transactions (id TEXT PRIMARY KEY, account_id_from TEXT, account_id_to TEXT, amount INTEGER, created_at DATETIME)")

	// Create clients
	client, err := entity.NewClient("any_name", "any_email")
	s.Nil(err)
	s.clientFrom = client

	client, err = entity.NewClient("any_name2", "any_email2")
	s.Nil(err)
	s.clientTo = client

	// Create accounts
	accountFrom := entity.NewAccount(s.clientFrom)
	accountFrom.Balance = 1000
	s.accountFrom = accountFrom

	accountTo := entity.NewAccount(s.clientTo)
	accountTo.Balance = 1000
	s.accountTo = accountTo

	s.TransactionDB = NewTransactionDB(db)
}

func (s *TransactionDBTestSuite) TearDownSuite() {
	defer s.db.Close()
	s.db.Exec("DROP TABLE clients")
	s.db.Exec("DROP TABLE accounts")
	s.db.Exec("DROP TABLE transactions")
}

func TestTransactionDBTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionDBTestSuite))
}

func (s *TransactionDBTestSuite) TestCreate() {
	transaction, err := entity.NewTransaction(s.accountFrom, s.accountTo, 100)
	s.Nil(err)

	err = s.TransactionDB.Create(transaction)
	s.Nil(err)
}
