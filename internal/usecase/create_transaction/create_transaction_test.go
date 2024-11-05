package create_transaction

import (
	"context"
	"testing"

	"github.com/mwives/microservices-fc-walletcore/internal/entity"
	"github.com/mwives/microservices-fc-walletcore/internal/event"
	"github.com/mwives/microservices-fc-walletcore/internal/usecase/mocks"
	"github.com/mwives/microservices-fc-walletcore/pkg/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateTransactionUseCase_Execute(t *testing.T) {
	clientAccountFrom, _ := entity.NewClient("any_name", "any_email")
	accountFrom := entity.NewAccount(clientAccountFrom)
	accountFrom.Credit(100)

	clientAccountTo, _ := entity.NewClient("any_name2", "any_email2")
	accountTo := entity.NewAccount(clientAccountTo)
	accountTo.Credit(100)

	uowMock := &mocks.UowMock{}
	uowMock.On("Do", mock.Anything, mock.Anything).Return(nil)

	dispatcherMock := events.NewEventDispatcher()
	event := event.NewTransactionCreatedEvent()
	ctx := context.Background()

	uc := NewCreateTransactionUseCase(uowMock, dispatcherMock, event)

	input := CreateTransactionInputDTO{
		AccountIDFrom: accountFrom.ID,
		AccountIDTo:   accountTo.ID,
		Amount:        10,
	}

	output, err := uc.Execute(ctx, input)
	assert.Nil(t, err)

	assert.NotNil(t, output)
	uowMock.AssertExpectations(t)
	uowMock.AssertNumberOfCalls(t, "Do", 1)
}
