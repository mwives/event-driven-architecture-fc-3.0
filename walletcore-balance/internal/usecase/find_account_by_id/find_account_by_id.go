package find_account_by_id

import "github.com/mwives/event-driven-architecture-fc-3.0/walletcore-balance/internal/gateway"

type FindAccountByIDInputDTO struct {
	AccountID string
}

type FindAccountByIDOutputDTO struct {
	AccountID string  `json:"account_id"`
	Balance   float64 `json:"balance"`
}

type FindAccountByIDUseCase struct {
	AccountGateway gateway.AccountGateway
}

func NewFindAccountByIDUseCase(accountGateway gateway.AccountGateway) *FindAccountByIDUseCase {
	return &FindAccountByIDUseCase{
		AccountGateway: accountGateway,
	}
}

func (uc *FindAccountByIDUseCase) Execute(input FindAccountByIDInputDTO) (*FindAccountByIDOutputDTO, error) {
	account, err := uc.AccountGateway.FindByID(input.AccountID)
	if err != nil {
		return nil, err
	}

	return &FindAccountByIDOutputDTO{
		AccountID: account.ID,
		Balance:   account.Balance,
	}, nil
}
