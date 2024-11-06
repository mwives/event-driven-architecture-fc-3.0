package create_account

import (
	"log"

	"github.com/mwives/event-driven-architecture-fc-3.0/walletcore/internal/entity"
	"github.com/mwives/event-driven-architecture-fc-3.0/walletcore/internal/gateway"
)

type CreateAccountInputDTO struct {
	ClientID string `json:"client_id"`
}

type CreateAccountOutputDTO struct {
	ID string
}

type CreateAccountUseCase struct {
	AcountGateway gateway.AccountGateway
	ClientGateway gateway.ClientGateway
}

func NewCreateAccountUseCase(
	accountGateway gateway.AccountGateway, clientGateway gateway.ClientGateway,
) *CreateAccountUseCase {
	return &CreateAccountUseCase{
		AcountGateway: accountGateway,
		ClientGateway: clientGateway,
	}
}

func (uc *CreateAccountUseCase) Execute(input CreateAccountInputDTO) (*CreateAccountOutputDTO, error) {
	client, err := uc.ClientGateway.FindByID(input.ClientID)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	account := entity.NewAccount(client)

	err = uc.AcountGateway.Create(account)
	if err != nil {
		return nil, err
	}

	return &CreateAccountOutputDTO{
		ID: account.ID,
	}, nil
}
