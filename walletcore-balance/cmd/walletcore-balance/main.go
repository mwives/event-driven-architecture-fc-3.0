package main

import (
	"database/sql"
	"fmt"
	"log"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mwives/event-driven-architecture-fc-3.0/walletcore-balance/internal/database"
	"github.com/mwives/event-driven-architecture-fc-3.0/walletcore-balance/internal/event/handler"
	"github.com/mwives/event-driven-architecture-fc-3.0/walletcore-balance/internal/usecase/find_account_by_id"
	"github.com/mwives/event-driven-architecture-fc-3.0/walletcore-balance/internal/usecase/sync_account_balance"
	"github.com/mwives/event-driven-architecture-fc-3.0/walletcore-balance/internal/web"
	"github.com/mwives/event-driven-architecture-fc-3.0/walletcore-balance/internal/web/webserver"
	"github.com/mwives/event-driven-architecture-fc-3.0/walletcore-balance/pkg/kafka"
)

func main() {
	db, err := sql.Open("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		"root", "root", "balance-db", "3306", "balance?parseTime=true",
	))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	configMap := ckafka.ConfigMap{
		"bootstrap.servers": "kafka:29092",
		"group.id":          "wallet",
	}
	topics := []string{"balances"}
	kafkaConsumer := kafka.NewConsumer(&configMap, topics)

	accountDB := database.NewAccountDB(db)
	createUpdateAccountUseCase := sync_account_balance.NewSyncAccountBalanceUseCase(accountDB)
	findAccountByIdUseCase := find_account_by_id.NewFindAccountByIDUseCase(accountDB)

	accountHandler := web.NewWebAccountHandler(*findAccountByIdUseCase)

	webServer := webserver.NewWebServer(":8080")
	webServer.AddHandler("/balances/{id}", accountHandler.FindAccountByID)
	go webServer.Start()

	msgChan := make(chan *ckafka.Message)

	go func() {
		if err := kafkaConsumer.Consume(msgChan); err != nil {
			log.Fatalf("Error while consuming: %v", err)
		}
	}()

	for msg := range msgChan {
		switch *msg.TopicPartition.Topic {
		case "balances":
			log.Printf("Message on %s: %s", *msg.TopicPartition.Topic, string(msg.Value))
			handler.ProcessBalanceUpdated(msg, createUpdateAccountUseCase)
		default:
			log.Printf("Unhandled topic: %s", *msg.TopicPartition.Topic)
		}
	}
}
