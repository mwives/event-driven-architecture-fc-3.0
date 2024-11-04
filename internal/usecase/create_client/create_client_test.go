package create_client

import (
	"testing"

	"github.com/mwives/microservices-fc-walletcore/internal/usecase/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateClientUseCase_Execute(t *testing.T) {
	clientGatewayMock := &mocks.ClientGatewayMock{}
	clientGatewayMock.On("Create", mock.Anything).Return(nil)

	uc := NewCreateClientUseCase(clientGatewayMock)

	input := &CreateClientInputDTO{
		Name:  "any_name",
		Email: "any_email",
	}

	output, err := uc.Execute(input)
	assert.Nil(t, err)

	assert.Equal(t, "any_name", output.Name)
	assert.Equal(t, "any_email", output.Email)
	assert.NotEmpty(t, output.ID)
	assert.NotZero(t, output.CreatedAt)
	assert.NotZero(t, output.UpdateAt)

	clientGatewayMock.AssertExpectations(t)
	clientGatewayMock.AssertNumberOfCalls(t, "Create", 1)
}
