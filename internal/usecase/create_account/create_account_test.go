package create_account

import (
	"testing"

	"github.com/mwives/microservices-fc-walletcore/internal/entity"
	"github.com/mwives/microservices-fc-walletcore/internal/usecase/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateAccountUseCase_Execute(t *testing.T) {
	client, _ := entity.NewClient("any_name", "any_email")
	clientGatewayMock := &mocks.ClientGatewayMock{}
	clientGatewayMock.On("FindByID", mock.Anything).Return(client, nil)

	accountGatewayMock := &mocks.AccountGatewayMock{}
	accountGatewayMock.On("Create", mock.Anything).Return(nil)

	uc := NewCreateAccountUseCase(accountGatewayMock, clientGatewayMock)

	input := CreateAccountInputDTO{
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
