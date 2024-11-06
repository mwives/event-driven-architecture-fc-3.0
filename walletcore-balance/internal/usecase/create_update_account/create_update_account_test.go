package create_update_account

import (
	"testing"
	"time"

	"github.com/mwives/event-driven-architecture-fc-3.0/walletcore-balance/internal/entity"
	"github.com/mwives/event-driven-architecture-fc-3.0/walletcore-balance/internal/usecase/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateUpdateAccountUseCase_CreateNewAccount(t *testing.T) {
	mockAccountGateway := new(mocks.MockAccountGateway)
	usecase := NewCreateUpdateAccountUseCase(mockAccountGateway)

	input := CreateUpdateAccountInputDTO{
		AccountID: "new_account",
		Balance:   1000,
	}

	mockAccountGateway.On("FindByID", input.AccountID).Return(&entity.Account{}, nil)
	mockAccountGateway.On("Create", mock.Anything).Return(nil)

	output, err := usecase.Execute(input)
	assert.Nil(t, err)

	assert.NotNil(t, output)
	assert.Equal(t, output.AccountID, input.AccountID)
	assert.Equal(t, output.Balance, input.Balance)

	mockAccountGateway.AssertExpectations(t)
	mockAccountGateway.AssertNumberOfCalls(t, "FindByID", 1)
	mockAccountGateway.AssertNumberOfCalls(t, "Create", 1)
	mockAccountGateway.AssertNumberOfCalls(t, "Update", 0)
}

func TestCreateUpdateAccountUseCase_UpdateAccount(t *testing.T) {
	mockAccountGateway := new(mocks.MockAccountGateway)
	usecase := NewCreateUpdateAccountUseCase(mockAccountGateway)

	existingAccount := &entity.Account{
		ID:        "existing_account",
		Balance:   2000,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	input := CreateUpdateAccountInputDTO{
		AccountID: existingAccount.ID,
		Balance:   1000,
	}

	// return account with current balance
	mockAccountGateway.On("FindByID", input.AccountID).Return(existingAccount, nil)

	existingAccount.Balance += 1000 // assuming update balance of 1000
	mockAccountGateway.On("Update", mock.Anything).Return(nil)

	output, err := usecase.Execute(input)
	assert.Nil(t, err)

	assert.NotNil(t, output)
	assert.Equal(t, output.AccountID, existingAccount.ID)
	assert.Equal(t, output.Balance, existingAccount.Balance)

	mockAccountGateway.AssertExpectations(t)
	mockAccountGateway.AssertNumberOfCalls(t, "FindByID", 1)
	mockAccountGateway.AssertNumberOfCalls(t, "Create", 0)
	mockAccountGateway.AssertNumberOfCalls(t, "Update", 1)
}
