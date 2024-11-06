package create_transaction

import (
	"context"

	"github.com/mwives/event-driven-architecture-fc-3.0/walletcore/internal/entity"
	"github.com/mwives/event-driven-architecture-fc-3.0/walletcore/internal/gateway"
	"github.com/mwives/event-driven-architecture-fc-3.0/walletcore/pkg/events"
	"github.com/mwives/event-driven-architecture-fc-3.0/walletcore/pkg/uow"
)

type CreateTransactionInputDTO struct {
	AccountIDFrom string  `json:"account_id_from"`
	AccountIDTo   string  `json:"account_id_to"`
	Amount        float64 `json:"amount"`
}

type CreateTransactionOutputDTO struct {
	ID            string
	AccountIDFrom string
	AccountIDTo   string
	Amount        float64
}

type UpdateBalanceOutputDTO struct {
	AccountFromID      string  `json:"account_id_from"`
	AccountToID        string  `json:"account_id_to"`
	AccountFromBalance float64 `json:"balance_account_id_from"`
	AccountToBalance   float64 `json:"balance_account_id_to"`
}

type CreateTransactionUseCase struct {
	Uow                uow.UowInterface
	EventDispatcher    events.EventDispatcherInterface
	TransactionCreated events.EventInterface
	BalanceUpdated     events.EventInterface
}

func NewCreateTransactionUseCase(
	uow uow.UowInterface,
	eventDispatcher events.EventDispatcherInterface,
	transactionCreated events.EventInterface,
	balanceUpdated events.EventInterface,
) *CreateTransactionUseCase {
	return &CreateTransactionUseCase{
		Uow:                uow,
		EventDispatcher:    eventDispatcher,
		TransactionCreated: transactionCreated,
		BalanceUpdated:     balanceUpdated,
	}
}

func (uc *CreateTransactionUseCase) Execute(ctx context.Context, input CreateTransactionInputDTO) (*CreateTransactionOutputDTO, error) {
	createTransactionOutput := &CreateTransactionOutputDTO{}
	updateBalanceOutput := &UpdateBalanceOutputDTO{}

	err := uc.Uow.Do(ctx, func(_ *uow.Uow) error {
		accountRepository := uc.getAccountRepository(ctx)
		transactionRepository := uc.getTransactionRepository(ctx)

		accountFrom, err := accountRepository.FindByID(input.AccountIDFrom)
		if err != nil {
			return err
		}
		accountTo, err := accountRepository.FindByID(input.AccountIDTo)
		if err != nil {
			return err
		}

		transaction, err := entity.NewTransaction(accountFrom, accountTo, input.Amount)
		if err != nil {
			return err
		}

		err = accountRepository.UpdateBalance(accountFrom)
		if err != nil {
			return err
		}
		err = accountRepository.UpdateBalance(accountTo)
		if err != nil {
			return err
		}

		err = transactionRepository.Create(transaction)
		if err != nil {
			return err
		}

		createTransactionOutput.ID = transaction.ID
		createTransactionOutput.AccountIDFrom = transaction.AccountFrom.ID
		createTransactionOutput.AccountIDTo = transaction.AccountTo.ID
		createTransactionOutput.Amount = transaction.Amount

		updateBalanceOutput.AccountFromID = accountFrom.ID
		updateBalanceOutput.AccountToID = accountTo.ID
		updateBalanceOutput.AccountFromBalance = accountFrom.Balance
		updateBalanceOutput.AccountToBalance = accountTo.Balance

		return nil
	})
	if err != nil {
		return nil, err
	}

	uc.TransactionCreated.SetPayload(createTransactionOutput)
	uc.EventDispatcher.Dispatch(uc.TransactionCreated)

	uc.BalanceUpdated.SetPayload(updateBalanceOutput)
	uc.EventDispatcher.Dispatch(uc.BalanceUpdated)

	return createTransactionOutput, nil
}

func (uc *CreateTransactionUseCase) getAccountRepository(ctx context.Context) gateway.AccountGateway {
	repo, err := uc.Uow.GetRepository(ctx, "AccountDB")
	if err != nil {
		panic(err)
	}
	return repo.(gateway.AccountGateway)
}

func (uc *CreateTransactionUseCase) getTransactionRepository(ctx context.Context) gateway.TransactionGateway {
	repo, err := uc.Uow.GetRepository(ctx, "TransactionDB")
	if err != nil {
		panic(err)
	}
	return repo.(gateway.TransactionGateway)
}
