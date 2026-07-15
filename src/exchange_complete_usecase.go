package main

import (
	"context"
	"fmt"
)

func (useCase *ExchangeUseCase) Complete(ctx context.Context, actorID, exchangeID int) (Exchange, error) {
	exec := useCase.db.Executor()

	exchange, err := useCase.exchanges.FindByID(ctx, exec, exchangeID)
	if err != nil {
		return Exchange{}, err
	}
	if actorID != exchange.OwnerID {
		return Exchange{}, ErrForbidden
	}
	if exchange.Status != ExchangeStatusAccepted {
		return Exchange{}, ErrExchangeInvalidTransition
	}

	service, err := useCase.services.FindByID(ctx, exec, exchange.ServiceID)
	if err != nil {
		return Exchange{}, err
	}

	var updated Exchange

	err = useCase.db.WithinTransaction(ctx, func(transactionExec dbExecutor) error {
		updated, err = useCase.exchanges.UpdateStatus(ctx, transactionExec, exchangeID, ExchangeStatusCompleted, exchangeUpdatedNow())
		if err != nil {
			return err
		}

		earn := creditTransactionForExchange(exchange.OwnerID, exchangeID, service.Credits, "earn")
		return useCase.creditTransactions.Create(ctx, transactionExec, earn)
	})
	if err != nil {
		return Exchange{}, fmt.Errorf("complete exchange: %w", err)
	}

	return updated, nil
}
