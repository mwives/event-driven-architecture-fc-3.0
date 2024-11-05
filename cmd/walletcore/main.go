package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mwives/microservices-fc-walletcore/internal/database"
	"github.com/mwives/microservices-fc-walletcore/internal/event"
	"github.com/mwives/microservices-fc-walletcore/internal/usecase/create_account"
	"github.com/mwives/microservices-fc-walletcore/internal/usecase/create_client"
	"github.com/mwives/microservices-fc-walletcore/internal/usecase/create_transaction"
	"github.com/mwives/microservices-fc-walletcore/internal/web"
	"github.com/mwives/microservices-fc-walletcore/internal/web/webserver"
	"github.com/mwives/microservices-fc-walletcore/pkg/events"
)

func main() {
	db, err := sql.Open("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		"root", "root", "localhost", "3306", "wallet",
	))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	eventDispatcher := events.NewEventDispatcher()
	transactionCreatedEvent := event.NewTransactionCreatedEvent()
	// eventDispatcher.Register("transaction_created", handler)

	clientDB := database.NewClientDB(db)
	accountDB := database.NewAccountDB(db)
	transactionDB := database.NewTransactionDB(db)

	createAccountUseCase := create_account.NewCreateAccountUseCase(accountDB, clientDB)
	createClientUseCase := create_client.NewCreateClientUseCase(clientDB)
	createTransactionUseCase := create_transaction.NewCreateTransactionUseCase(
		transactionDB, accountDB, eventDispatcher, transactionCreatedEvent,
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
