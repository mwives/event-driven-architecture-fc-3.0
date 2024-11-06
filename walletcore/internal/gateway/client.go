package gateway

import "github.com/mwives/event-driven-architecture-fc-3.0/walletcore/internal/entity"

type ClientGateway interface {
	Create(client *entity.Client) error
	FindByID(id string) (*entity.Client, error)
}
