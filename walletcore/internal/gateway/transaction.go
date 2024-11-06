package gateway

import "github.com/mwives/event-driven-architecture-fc-3.0/walletcore/internal/entity"

type TransactionGateway interface {
	Create(transaction *entity.Transaction) error
}
