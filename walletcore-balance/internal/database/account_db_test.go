package database

import (
	"database/sql"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mwives/event-driven-architecture-fc-3.0/walletcore-balance/internal/entity"
	"github.com/stretchr/testify/suite"
)

type AccountDbTestSuite struct {
	suite.Suite
	db        *sql.DB
	accountDB *AccountDB
}

func (s *AccountDbTestSuite) SetupSuite() {
	dsn := "root:root@tcp(localhost:3307)/wallet_balance_test?parseTime=true"
	db, err := sql.Open("mysql", dsn)
	s.Nil(err)
	s.db = db

	db.Exec(
		`CREATE TABLE
			accounts (
				id VARCHAR(255) PRIMARY KEY,
				balance INT,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
			)`,
	)

	s.accountDB = NewAccountDB(db)
}

func (s *AccountDbTestSuite) TearDownSuite() {
	defer s.db.Close()
	s.db.Exec("DROP TABLE accounts")
}

func TestAccountDbTestSuite(t *testing.T) {
	suite.Run(t, new(AccountDbTestSuite))
}

func (s *AccountDbTestSuite) TestCreate() {
	account := entity.NewAccount()
	err := s.accountDB.Create(account)
	s.Nil(err)
}

func (s *AccountDbTestSuite) TestFindByID() {
	account := entity.NewAccount()
	err := s.accountDB.Create(account)
	s.Nil(err)

	foundAccount, err := s.accountDB.FindByID(account.ID)
	s.Nil(err)
	s.Equal(account.ID, foundAccount.ID)
	s.Equal(account.Balance, foundAccount.Balance)

	// Normalize both timestamps to UTC for comparison
	accountCreatedAtUTC := account.CreatedAt.UTC()
	foundAccountCreatedAtUTC := foundAccount.CreatedAt.UTC()
	accountUpdatedAtUTC := account.UpdatedAt.UTC()
	foundAccountUpdatedAtUTC := foundAccount.UpdatedAt.UTC()

	// Compare timestamps with a tolerance of 1 second (milliseconds are not stored in MySQL)
	s.WithinDuration(accountCreatedAtUTC, foundAccountCreatedAtUTC, time.Second)
	s.WithinDuration(accountUpdatedAtUTC, foundAccountUpdatedAtUTC, time.Second)
}

func (s *AccountDbTestSuite) TestUpdate() {
	account := entity.NewAccount()
	err := s.accountDB.Create(account)
	s.Nil(err)

	account.Balance = 100
	err = s.accountDB.Update(account)
	s.Nil(err)

	foundAccount, err := s.accountDB.FindByID(account.ID)
	s.Nil(err)
	s.Equal(account.Balance, foundAccount.Balance)
}
