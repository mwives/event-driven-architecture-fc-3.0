package create_transaction

import (
	"context"
	"testing"

	"github.com/mwives/event-driven-architecture-fc-3.0/walletcore/internal/entity"
	"github.com/mwives/event-driven-architecture-fc-3.0/walletcore/internal/event"
	"github.com/mwives/event-driven-architecture-fc-3.0/walletcore/internal/usecase/mocks"
	"github.com/mwives/event-driven-architecture-fc-3.0/walletcore/pkg/events"
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
	transactionCreatedEvent := event.NewTransactionCreatedEvent()
	balanceUpdatedEvent := event.NewBalanceUpdatedEvent()
	ctx := context.Background()

	uc := NewCreateTransactionUseCase(uowMock, dispatcherMock, transactionCreatedEvent, balanceUpdatedEvent)

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
