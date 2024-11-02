package gateway

import "github.com/mwives/microservices-fc-walletcore/internal/entity"

type ClientGateway interface {
	Create(client *entity.Client) error
	FindByID(id string) (*entity.Client, error)
}
