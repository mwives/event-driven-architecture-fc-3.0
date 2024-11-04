package create_account

import (
	"testing"

	"github.com/mwives/microservices-fc-walletcore/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type ClientGatewayMock struct {
	mock.Mock
}

func (m *ClientGatewayMock) Create(client *entity.Client) error {
	args := m.Called(client)
	return args.Error(0)
}

func (m *ClientGatewayMock) FindByID(id string) (*entity.Client, error) {
	args := m.Called(id)
	return args.Get(0).(*entity.Client), args.Error(1)
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

func TestCreateAccountUseCase_Execute(t *testing.T) {
	client, _ := entity.NewClient("any_name", "any_email")
	clientGatewayMock := &ClientGatewayMock{}
	clientGatewayMock.On("FindByID", mock.Anything).Return(client, nil)

	accountGatewayMock := &AccountGatewayMock{}
	accountGatewayMock.On("Create", mock.Anything).Return(nil)

	uc := NewCreateAccountUseCase(accountGatewayMock, clientGatewayMock)

	input := &CreateAccountInputDTO{
		ClientID: client.ID,
	}

	output, err := uc.Execute(input)
	assert.Nil(t, err)

	assert.NotEmpty(t, output.ID)

	clientGatewayMock.AssertExpectations(t)
	clientGatewayMock.AssertNumberOfCalls(t, "FindByID", 1)

	accountGatewayMock.AssertExpectations(t)
	accountGatewayMock.AssertNumberOfCalls(t, "Create", 1)
}
