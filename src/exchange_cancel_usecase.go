package main

import (
	"context"
	"fmt"
)

func (useCase *ExchangeUseCase) Cancel(ctx context.Context, actorID, exchangeID int) (Exchange, error) {
	exec := useCase.db.Executor()

	exchange, err := useCase.exchanges.FindByID(ctx, exec, exchangeID)
	if err != nil {
		return Exchange{}, err
	}
	if actorID != exchange.RequesterID && actorID != exchange.OwnerID {
		return Exchange{}, ErrForbidden
	}
	if exchange.Status != ExchangeStatusPending && exchange.Status != ExchangeStatusAccepted {
		return Exchange{}, ErrExchangeInvalidTransition
	}

	if exchange.Status == ExchangeStatusPending {
		updated, err := useCase.exchanges.UpdateStatus(ctx, exec, exchangeID, ExchangeStatusCancelled, exchangeUpdatedNow())
		if err != nil {
			return Exchange{}, fmt.Errorf("cancel exchange: %w", err)
		}

		return updated, nil
	}

	service, err := useCase.services.FindByID(ctx, exec, exchange.ServiceID)
	if err != nil {
		return Exchange{}, err
	}

	var updated Exchange

	err = useCase.db.WithinTransaction(ctx, func(transactionExec dbExecutor) error {
		updated, err = useCase.exchanges.UpdateStatus(ctx, transactionExec, exchangeID, ExchangeStatusCancelled, exchangeUpdatedNow())
		if err != nil {
			return err
		}

		refund := creditTransactionForExchange(exchange.RequesterID, exchangeID, service.Credits, "refund")
		return useCase.creditTransactions.Create(ctx, transactionExec, refund)
	})
	if err != nil {
		return Exchange{}, fmt.Errorf("cancel exchange: %w", err)
	}

	return updated, nil
}
