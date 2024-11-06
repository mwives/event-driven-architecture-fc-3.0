package find_account_by_id

import (
	"testing"

	"github.com/mwives/event-driven-architecture-fc-3.0/walletcore-balance/internal/entity"
	"github.com/mwives/event-driven-architecture-fc-3.0/walletcore-balance/internal/usecase/mocks"
	"github.com/stretchr/testify/assert"
)

func TestFindAccountByIdUseCase_Execute(t *testing.T) {
	mockAccountGateway := new(mocks.MockAccountGateway)
	usecase := NewFindAccountByIDUseCase(mockAccountGateway)

	account := &entity.Account{
		ID:      "account_id",
		Balance: 1000,
	}

	input := FindAccountByIDInputDTO{
		AccountID: account.ID,
	}

	mockAccountGateway.On("FindByID", input.AccountID).Return(account, nil)

	output, err := usecase.Execute(input)
	assert.Nil(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, input.AccountID, output.AccountID)

	mockAccountGateway.AssertExpectations(t)
	mockAccountGateway.AssertNumberOfCalls(t, "FindByID", 1)
}
