package gateway

import "github.com/mwives/microservices-fc-walletcore/internal/entity"

type TransactionGateway interface {
	Create(transaction *entity.Transaction) error
}
