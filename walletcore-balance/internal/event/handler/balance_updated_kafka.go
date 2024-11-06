package handler

import (
	"encoding/json"
	"log"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/mwives/event-driven-architecture-fc-3.0/walletcore-balance/internal/usecase/sync_account_balance"
)

type messagePayloadDTO struct {
	Name    string
	Payload sync_account_balance.SyncAccountBalanceInputDTO
}

func ProcessBalanceUpdated(msg *ckafka.Message, useCase *sync_account_balance.SyncAccountBalanceUseCase) {
	var inputDTO messagePayloadDTO

	if err := json.Unmarshal(msg.Value, &inputDTO); err != nil {
		log.Printf("Failed to unmarshal balance update message: %v", err)
		return
	}

	output, err := useCase.Execute(inputDTO.Payload)
	if err != nil {
		log.Printf("Error executing balance update use case: %v", err)
		return
	}

	log.Printf("Processed balance update: %+v", output)
}
