package gateway

import "github.com/mwives/event-driven-architecture-fc-3.0/walletcore/internal/entity"

type AccountGateway interface {
	Create(account *entity.Account) error
	FindByID(id string) (*entity.Account, error)
	UpdateBalance(account *entity.Account) error
}
