package sync_account_balance

import (
	"testing"
	"time"

	"github.com/mwives/event-driven-architecture-fc-3.0/walletcore-balance/internal/entity"
	"github.com/mwives/event-driven-architecture-fc-3.0/walletcore-balance/internal/usecase/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSyncAccountBalanceUseCase_CreateNewAccounts(t *testing.T) {
	mockAccountGateway := new(mocks.MockAccountGateway)
	usecase := NewSyncAccountBalanceUseCase(mockAccountGateway)

	input := SyncAccountBalanceInputDTO{
		AccountFromID:      "new_account_from",
		AccountToID:        "new_account_to",
		AccountFromBalance: 1000,
		AccountToBalance:   2000,
	}

	// Mock creating accounts
	mockAccountGateway.On("FindByID", input.AccountFromID).Return(&entity.Account{}, nil)
	mockAccountGateway.On("FindByID", input.AccountToID).Return(&entity.Account{}, nil)
	mockAccountGateway.On("Create", mock.Anything).Return(nil)

	output, err := usecase.Execute(input)
	assert.Nil(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, input.AccountFromID, output.AccountFromID)
	assert.Equal(t, input.AccountToID, output.AccountToID)
	assert.Equal(t, input.AccountFromBalance, output.AccountFromBalance)
	assert.Equal(t, input.AccountToBalance, output.AccountToBalance)

	mockAccountGateway.AssertExpectations(t)
	mockAccountGateway.AssertNumberOfCalls(t, "FindByID", 2)
	mockAccountGateway.AssertNumberOfCalls(t, "Create", 2)
	mockAccountGateway.AssertNumberOfCalls(t, "Update", 0)
}

func TestSyncAccountBalanceUseCase_UpdateExistingAccounts(t *testing.T) {
	mockAccountGateway := new(mocks.MockAccountGateway)
	usecase := NewSyncAccountBalanceUseCase(mockAccountGateway)

	existingAccountFrom := &entity.Account{
		ID:        "existing_account_from",
		Balance:   1000,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	existingAccountTo := &entity.Account{
		ID:        "existing_account_to",
		Balance:   2000,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	input := SyncAccountBalanceInputDTO{
		AccountFromID:      existingAccountFrom.ID,
		AccountToID:        existingAccountTo.ID,
		AccountFromBalance: 1500,
		AccountToBalance:   2500,
	}

	// Return existing accounts
	mockAccountGateway.On("FindByID", input.AccountFromID).Return(existingAccountFrom, nil)
	mockAccountGateway.On("FindByID", input.AccountToID).Return(existingAccountTo, nil)

	// Update balances
	mockAccountGateway.On("Update", mock.Anything).Return(nil)

	output, err := usecase.Execute(input)
	assert.Nil(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, input.AccountFromID, output.AccountFromID)
	assert.Equal(t, input.AccountToID, output.AccountToID)
	assert.Equal(t, input.AccountFromBalance, output.AccountFromBalance)
	assert.Equal(t, input.AccountToBalance, output.AccountToBalance)

	mockAccountGateway.AssertExpectations(t)
	mockAccountGateway.AssertNumberOfCalls(t, "FindByID", 2)
	mockAccountGateway.AssertNumberOfCalls(t, "Create", 0)
	mockAccountGateway.AssertNumberOfCalls(t, "Update", 2)
}
