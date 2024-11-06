package main

import (
	"context"
	"database/sql"
	"fmt"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mwives/event-driven-architecture-fc-3.0/walletcore/internal/database"
	"github.com/mwives/event-driven-architecture-fc-3.0/walletcore/internal/event"
	"github.com/mwives/event-driven-architecture-fc-3.0/walletcore/internal/event/handler"
	"github.com/mwives/event-driven-architecture-fc-3.0/walletcore/internal/usecase/create_account"
	"github.com/mwives/event-driven-architecture-fc-3.0/walletcore/internal/usecase/create_client"
	"github.com/mwives/event-driven-architecture-fc-3.0/walletcore/internal/usecase/create_transaction"
	"github.com/mwives/event-driven-architecture-fc-3.0/walletcore/internal/web"
	"github.com/mwives/event-driven-architecture-fc-3.0/walletcore/internal/web/webserver"
	"github.com/mwives/event-driven-architecture-fc-3.0/walletcore/pkg/events"
	"github.com/mwives/event-driven-architecture-fc-3.0/walletcore/pkg/kafka"
	"github.com/mwives/event-driven-architecture-fc-3.0/walletcore/pkg/uow"
)

func main() {
	db, err := sql.Open("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		"root", "root", "wallet-db", "3306", "wallet",
	))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	configMap := ckafka.ConfigMap{
		"bootstrap.servers": "kafka:29092",
		"group.id":          "wallet",
	}
	kafkaProducer := kafka.NewKafkaProducer(&configMap)

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("transaction_created", handler.NewTransactionCreatedKafkaHandler(kafkaProducer))
	eventDispatcher.Register("balance_updated", handler.NewBalanceUpdatedKafkaHandler(kafkaProducer))
	transactionCreatedEvent := event.NewTransactionCreatedEvent()
	balanceUpdatedEvent := event.NewBalanceUpdatedEvent()

	clientDB := database.NewClientDB(db)
	accountDB := database.NewAccountDB(db)

	ctx := context.Background()
	uow := uow.NewUow(ctx, db)

	uow.Register("AccountDB", func(tx *sql.Tx) interface{} {
		return database.NewAccountDB(db)
	})
	uow.Register("TransactionDB", func(tx *sql.Tx) interface{} {
		return database.NewTransactionDB(db)
	})

	createAccountUseCase := create_account.NewCreateAccountUseCase(accountDB, clientDB)
	createClientUseCase := create_client.NewCreateClientUseCase(clientDB)
	createTransactionUseCase := create_transaction.NewCreateTransactionUseCase(
		uow, eventDispatcher, transactionCreatedEvent, balanceUpdatedEvent,
	)

	accountHandler := web.NewWebAccountHandler(*createAccountUseCase)
	clientHandler := web.NewWebClientHandler(*createClientUseCase)
	transactionHandler := web.NewWebTransactionHandler(*createTransactionUseCase)

	webServer := webserver.NewWebServer(":8080")

	webServer.AddHandler("/accounts", accountHandler.CreateAccount)
	webServer.AddHandler("/clients", clientHandler.CreateClient)
	webServer.AddHandler("/transactions", transactionHandler.CreateTransaction)

	webServer.Start()
}
