package gateway

import "github.com/mwives/microservices-fc-walletcore/internal/entity"

type AccountGateway interface {
	Create(account *entity.Account) error
	FindByID(id string) (*entity.Account, error)
	UpdateBalance(account *entity.Account) error
}
