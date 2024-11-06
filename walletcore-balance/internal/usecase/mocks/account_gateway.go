package mocks

import (
	"github.com/mwives/event-driven-architecture-fc-3.0/walletcore-balance/internal/entity"
	"github.com/stretchr/testify/mock"
)

type MockAccountGateway struct {
	mock.Mock
}

func (m *MockAccountGateway) FindByID(id string) (*entity.Account, error) {
	args := m.Called(id)
	return args.Get(0).(*entity.Account), args.Error(1)
}

func (m *MockAccountGateway) Create(account *entity.Account) error {
	args := m.Called(account)
	return args.Error(0)
}

func (m *MockAccountGateway) Update(account *entity.Account) error {
	args := m.Called(account)
	return args.Error(0)

}
