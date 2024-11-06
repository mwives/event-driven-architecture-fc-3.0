package gateway

import "github.com/mwives/event-driven-architecture-fc-3.0/walletcore-balance/internal/entity"

type AccountGateway interface {
	FindByID(id string) (*entity.Account, error)
	Create(account *entity.Account) error
	Update(account *entity.Account) error
}
