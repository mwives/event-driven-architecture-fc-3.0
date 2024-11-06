package mocks

import (
	"context"

	"github.com/mwives/event-driven-architecture-fc-3.0/walletcore/pkg/uow"
	"github.com/stretchr/testify/mock"
)

type UowMock struct {
	mock.Mock
}

func (m *UowMock) Register(name string, fc uow.RepositoryFactory) {
	m.Called(name, fc)
}

func (m *UowMock) Unregister(name string) {
	m.Called(name)
}

func (m *UowMock) GetRepository(ctx context.Context, name string) (interface{}, error) {
	args := m.Called(name)
	return args.Get(0), args.Error(1)
}

func (m *UowMock) Do(ctx context.Context, work func(*uow.Uow) error) error {
	args := m.Called(work)
	return args.Error(0)
}

func (m *UowMock) CommitOrRollback() error {
	args := m.Called()
	return args.Error(0)
}

func (m *UowMock) Rollback() error {
	args := m.Called()
	return args.Error(0)
}
