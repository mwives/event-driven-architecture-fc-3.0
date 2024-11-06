package sync_account_balance

import (
	"time"

	"github.com/mwives/event-driven-architecture-fc-3.0/walletcore-balance/internal/entity"
	"github.com/mwives/event-driven-architecture-fc-3.0/walletcore-balance/internal/gateway"
)

type SyncAccountBalanceInputDTO struct {
	AccountFromID      string  `json:"account_id_from"`
	AccountToID        string  `json:"account_id_to"`
	AccountFromBalance float64 `json:"account_from_balance"`
	AccountToBalance   float64 `json:"account_to_balance"`
}

type SyncAccountBalanceOutputDTO struct {
	AccountFromID      string  `json:"account_id_from"`
	AccountToID        string  `json:"account_id_to"`
	AccountFromBalance float64 `json:"account_from_balance"`
	AccountToBalance   float64 `json:"account_to_balance"`
}

type SyncAccountBalanceUseCase struct {
	AccountGateway gateway.AccountGateway
}

func NewSyncAccountBalanceUseCase(accountGateway gateway.AccountGateway) *SyncAccountBalanceUseCase {
	return &SyncAccountBalanceUseCase{
		AccountGateway: accountGateway,
	}
}

func (uc *SyncAccountBalanceUseCase) Execute(input SyncAccountBalanceInputDTO) (*SyncAccountBalanceOutputDTO, error) {
	outputAccountFrom, err := uc.updateBalance(input.AccountFromID, input.AccountFromBalance)
	if err != nil {
		return nil, err
	}

	outputAccountTo, err := uc.updateBalance(input.AccountToID, input.AccountToBalance)
	if err != nil {
		return nil, err
	}

	return &SyncAccountBalanceOutputDTO{
		AccountFromID:      outputAccountFrom.AccountID,
		AccountToID:        outputAccountTo.AccountID,
		AccountFromBalance: outputAccountFrom.Balance,
		AccountToBalance:   outputAccountTo.Balance,
	}, nil
}

type updateBalanceOutputDTO struct {
	AccountID string
	Balance   float64
}

func (uc *SyncAccountBalanceUseCase) updateBalance(id string, balance float64) (updateBalanceOutputDTO, error) {
	account, err := uc.AccountGateway.FindByID(id)
	if err != nil {
		return updateBalanceOutputDTO{}, err
	}

	if account.ID == "" {
		account = &entity.Account{
			ID:        id,
			Balance:   balance,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err = uc.AccountGateway.Create(account)
		if err != nil {
			return updateBalanceOutputDTO{}, err
		}

		return updateBalanceOutputDTO{
			AccountID: account.ID,
			Balance:   account.Balance,
		}, nil
	}

	account.ID = id
	account.Balance = balance
	account.UpdatedAt = time.Now()

	err = uc.AccountGateway.Update(account)
	if err != nil {
		return updateBalanceOutputDTO{}, err
	}

	return updateBalanceOutputDTO{
		AccountID: account.ID,
		Balance:   account.Balance,
	}, nil
}
