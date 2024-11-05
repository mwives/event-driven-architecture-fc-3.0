package create_transaction

import (
	"testing"

	"github.com/mwives/microservices-fc-walletcore/internal/entity"
	"github.com/mwives/microservices-fc-walletcore/internal/event"
	"github.com/mwives/microservices-fc-walletcore/internal/usecase/mocks"
	"github.com/mwives/microservices-fc-walletcore/pkg/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateTransactionUseCase_Execute(t *testing.T) {
	transactionGatewayMock := &mocks.TransactionGatewayMock{}
	transactionGatewayMock.On("Create", mock.Anything).Return(nil)

	clientAccountFrom, _ := entity.NewClient("any_name", "any_email")
	accountFrom := entity.NewAccount(clientAccountFrom)
	accountFrom.Credit(100)

	clientAccountTo, _ := entity.NewClient("any_name2", "any_email2")
	accountTo := entity.NewAccount(clientAccountTo)
	accountTo.Credit(100)

	accountGatewayMock := &mocks.AccountGatewayMock{}
	accountGatewayMock.On("FindByID", mock.Anything).Return(accountFrom, nil)
	accountGatewayMock.On("FindByID", mock.Anything).Return(accountTo, nil)

	dispatcherMock := events.NewEventDispatcher()
	event := event.NewTransactionCreatedEvent()

	uc := NewCreateTransactionUseCase(transactionGatewayMock, accountGatewayMock, dispatcherMock, event)

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
