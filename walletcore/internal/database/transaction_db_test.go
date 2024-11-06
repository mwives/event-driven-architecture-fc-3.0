package database

import (
	"database/sql"
	"testing"

	"github.com/mwives/event-driven-architecture-fc-3.0/walletcore/internal/entity"
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

	db.Exec("CREATE TABLE clients (id VARCHAR(255) PRIMARY KEY, name VARCHAR(255), email VARCHAR(255), created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP)")
	db.Exec("CREATE TABLE accounts (id VARCHAR(255) PRIMARY KEY, client_id VARCHAR(255), balance INT, created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP)")
	db.Exec("CREATE TABLE transactions (id VARCHAR(255) PRIMARY KEY, account_id_from VARCHAR(255), account_id_to VARCHAR(255), amount INT, created_at TIMESTAMP)")

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
