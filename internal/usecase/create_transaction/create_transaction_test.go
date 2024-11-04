package create_transaction

import (
	"testing"

	"github.com/mwives/microservices-fc-walletcore/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type TransactionGatewayMock struct {
	mock.Mock
}

func (m *TransactionGatewayMock) Create(transaction *entity.Transaction) error {
	args := m.Called(transaction)
	return args.Error(0)
}

type AccountGatewayMock struct {
	mock.Mock
}

func (m *AccountGatewayMock) Create(account *entity.Account) error {
	args := m.Called(account)
	return args.Error(0)
}

func (m *AccountGatewayMock) FindByID(id string) (*entity.Account, error) {
	args := m.Called(id)
	return args.Get(0).(*entity.Account), args.Error(1)
}

func TestCreateTransactionUseCase_Execute(t *testing.T) {
	transactionGatewayMock := &TransactionGatewayMock{}
	transactionGatewayMock.On("Create", mock.Anything).Return(nil)

	clientAccountFrom, _ := entity.NewClient("any_name", "any_email")
	accountFrom := entity.NewAccount(clientAccountFrom)
	accountFrom.Credit(100)

	clientAccountTo, _ := entity.NewClient("any_name2", "any_email2")
	accountTo := entity.NewAccount(clientAccountTo)
	accountTo.Credit(100)

	accountGatewayMock := &AccountGatewayMock{}
	accountGatewayMock.On("FindByID", mock.Anything).Return(accountFrom, nil)
	accountGatewayMock.On("FindByID", mock.Anything).Return(accountTo, nil)

	uc := NewCreateTransactionUseCase(transactionGatewayMock, accountGatewayMock)

	input := CreateTransactionInputDTO{
		AccountIDFrom: accountFrom.ID,
		AccountIDTo:   accountTo.ID,
		Amount:        10,
	}

	output, err := uc.Execute(input)
	assert.Nil(t, err)

	assert.NotEmpty(t, output.ID)

	accountGatewayMock.AssertExpectations(t)
	accountGatewayMock.AssertNumberOfCalls(t, "FindByID", 2)

	transactionGatewayMock.AssertExpectations(t)
	transactionGatewayMock.AssertNumberOfCalls(t, "Create", 1)
}
