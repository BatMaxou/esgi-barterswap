package main

import (
	"context"
	"fmt"
)

func (useCase *ExchangeUseCase) Accept(ctx context.Context, actorID, exchangeID int) (Exchange, error) {
	exec := useCase.db.Executor()

	exchange, err := useCase.exchanges.FindByID(ctx, exec, exchangeID)
	if err != nil {
		return Exchange{}, err
	}
	if actorID != exchange.OwnerID {
		return Exchange{}, ErrForbidden
	}
	if exchange.Status != ExchangeStatusPending {
		return Exchange{}, ErrExchangeInvalidTransition
	}

	service, err := useCase.services.FindByID(ctx, exec, exchange.ServiceID)
	if err != nil {
		return Exchange{}, err
	}

	var updated Exchange

	err = useCase.db.WithinTransaction(ctx, func(transactionExec dbExecutor) error {
		balance, err := useCase.creditTransactions.BalanceByUserID(ctx, transactionExec, exchange.RequesterID)
		if err != nil {
			return err
		}
		if balance < service.Credits {
			return ErrExchangeInsufficientCredits
		}

		updated, err = useCase.exchanges.UpdateStatus(ctx, transactionExec, exchangeID, ExchangeStatusAccepted, exchangeUpdatedNow())
		if err != nil {
			return err
		}

		spend := creditTransactionForExchange(exchange.RequesterID, exchangeID, -service.Credits, "spend")
		return useCase.creditTransactions.Create(ctx, transactionExec, spend)
	})
	if err != nil {
		return Exchange{}, fmt.Errorf("accept exchange: %w", err)
	}

	return updated, nil
}
