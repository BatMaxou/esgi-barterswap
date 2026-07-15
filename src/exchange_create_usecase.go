package main

import (
	"context"
	"fmt"
)

func (useCase *ExchangeUseCase) Create(ctx context.Context, requesterID, serviceID int) (Exchange, error) {
	if serviceID <= 0 {
		return Exchange{}, ErrExchangeServiceIDInvalid
	}

	exec := useCase.db.Executor()

	service, err := useCase.services.FindByID(ctx, exec, serviceID)
	if err != nil {
		return Exchange{}, err
	}
	if !service.Active {
		return Exchange{}, ErrExchangeServiceInactive
	}
	if requesterID == service.ProviderID {
		return Exchange{}, ErrExchangeSelfRequest
	}

	exchange, err := NewExchange(serviceID, requesterID, service.ProviderID)
	if err != nil {
		return Exchange{}, err
	}

	err = useCase.db.WithinTransaction(ctx, func(transactionExec dbExecutor) error {
		hasActive, err := useCase.exchanges.HasActiveForService(ctx, transactionExec, serviceID)
		if err != nil {
			return err
		}
		if hasActive {
			return ErrExchangeServiceUnavailable
		}

		balance, err := useCase.creditTransactions.BalanceByUserID(ctx, transactionExec, requesterID)
		if err != nil {
			return err
		}
		if balance < service.Credits {
			return ErrExchangeInsufficientCredits
		}

		created, err := useCase.exchanges.Create(ctx, transactionExec, exchange)
		if err != nil {
			return err
		}
		exchange = created

		return nil
	})
	if err != nil {
		return Exchange{}, fmt.Errorf("create exchange: %w", err)
	}

	return exchange, nil
}
