package handler

import (
	"encoding/json"
	"log"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/mwives/event-driven-architecture-fc-3.0/walletcore-balance/internal/usecase/sync_account_balance"
)

func ProcessBalanceUpdated(msg *ckafka.Message, useCase *sync_account_balance.SyncAccountBalanceUseCase) {
	var inputDTO sync_account_balance.SyncAccountBalanceInputDTO

	if err := json.Unmarshal(msg.Value, &inputDTO); err != nil {
		log.Printf("Failed to unmarshal balance update message: %v", err)
		return
	}

	output, err := useCase.Execute(inputDTO)
	if err != nil {
		log.Printf("Error executing balance update use case: %v", err)
		return
	}

	log.Printf("Processed balance update: %+v", output)
}
