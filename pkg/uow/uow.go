package uow

import (
	"context"
	"database/sql"
	"fmt"
)

type RepositoryFactory func(tx *sql.Tx) interface{}

type UowInterface interface {
	Register(name string, repoFactory RepositoryFactory)
	Unregister(name string)
	GetRepository(ctx context.Context, name string) (interface{}, error)
	Do(ctx context.Context, work func(uow *Uow) error) error
	CommitOrRollback() error
	Rollback() error
}

type Uow struct {
	Db           *sql.DB
	Tx           *sql.Tx
	Repositories map[string]RepositoryFactory
}

func NewUow(ctx context.Context, db *sql.DB) *Uow {
	return &Uow{
		Db:           db,
		Repositories: make(map[string]RepositoryFactory),
	}
}

func (u *Uow) Register(name string, repoFactory RepositoryFactory) {
	u.Repositories[name] = repoFactory
}

func (u *Uow) Unregister(name string) {
	delete(u.Repositories, name)
}

func (u *Uow) GetRepository(ctx context.Context, name string) (interface{}, error) {
	if u.Tx == nil {
		tx, err := u.Db.BeginTx(ctx, nil)
		if err != nil {
			return nil, err
		}
		u.Tx = tx
	}
	repo := u.Repositories[name](u.Tx)
	return repo, nil
}

func (u *Uow) Do(ctx context.Context, work func(uow *Uow) error) error {
	if u.Tx != nil {
		return fmt.Errorf("transaction already started")
	}

	tx, err := u.Db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	u.Tx = tx

	err = work(u)
	if err != nil {
		errRb := u.Rollback()
		if errRb != nil {
			return fmt.Errorf("rollback error: %s, work error: %s", errRb, err)
		}
		return err
	}

	return u.CommitOrRollback()
}

func (u *Uow) CommitOrRollback() error {
	if u.Tx == nil {
		return fmt.Errorf("no transaction to commit")
	}

	err := u.Tx.Commit()
	if err != nil {
		errRb := u.Rollback()
		if errRb != nil {
			return fmt.Errorf("rollback error: %s, commit error: %s", errRb, err)
		}
		return err
	}

	u.Tx = nil
	return nil
}

func (u *Uow) Rollback() error {
	if u.Tx == nil {
		return fmt.Errorf("no transaction to rollback")
	}

	err := u.Tx.Rollback()
	if err != nil {
		return err
	}

	u.Tx = nil
	return nil
}
