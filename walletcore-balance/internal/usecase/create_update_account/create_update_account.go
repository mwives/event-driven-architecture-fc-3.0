package create_update_account

import (
	"fmt"
	"time"

	"github.com/mwives/event-driven-architecture-fc-3.0/walletcore-balance/internal/entity"
	"github.com/mwives/event-driven-architecture-fc-3.0/walletcore-balance/internal/gateway"
)

type CreateUpdateAccountInputDTO struct {
	AccountID string  `json:"account_id"`
	Balance   float64 `json:"balance"`
}

type CreateUpdateAccountOutputDTO struct {
	AccountID string  `json:"account_id"`
	Balance   float64 `json:"balance"`
}

type CreateUpdateAccountUseCase struct {
	AccountGateway gateway.AccountGateway
}

func NewCreateUpdateAccountUseCase(accountGateway gateway.AccountGateway) *CreateUpdateAccountUseCase {
	return &CreateUpdateAccountUseCase{
		AccountGateway: accountGateway,
	}
}

func (uc *CreateUpdateAccountUseCase) Execute(input CreateUpdateAccountInputDTO) (*CreateUpdateAccountOutputDTO, error) {
	account, err := uc.AccountGateway.FindByID(input.AccountID)
	if err != nil {
		return nil, err
	}

	fmt.Printf("%v\n", account == nil)

	if account == nil {
		account = &entity.Account{
			ID:        input.AccountID,
			Balance:   input.Balance,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err = uc.AccountGateway.Create(account)
		if err != nil {
			return nil, err
		}

		return &CreateUpdateAccountOutputDTO{
			AccountID: account.ID,
			Balance:   account.Balance,
		}, nil
	}

	account.ID = input.AccountID
	account.Balance = input.Balance
	account.UpdatedAt = time.Now()

	err = uc.AccountGateway.Update(account)
	if err != nil {
		return nil, err
	}

	return &CreateUpdateAccountOutputDTO{
		AccountID: account.ID,
		Balance:   account.Balance,
	}, nil
}
