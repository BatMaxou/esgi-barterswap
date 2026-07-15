package main

import (
	"context"
	"fmt"
)

func (useCase *ExchangeUseCase) Reject(ctx context.Context, actorID, exchangeID int) (Exchange, error) {
	exchange, err := useCase.exchanges.FindByID(ctx, useCase.db.Executor(), exchangeID)
	if err != nil {
		return Exchange{}, err
	}
	if actorID != exchange.OwnerID {
		return Exchange{}, ErrForbidden
	}
	if exchange.Status != ExchangeStatusPending {
		return Exchange{}, ErrExchangeInvalidTransition
	}

	updated, err := useCase.exchanges.UpdateStatus(ctx, useCase.db.Executor(), exchangeID, ExchangeStatusRejected, exchangeUpdatedNow())
	if err != nil {
		return Exchange{}, fmt.Errorf("reject exchange: %w", err)
	}

	return updated, nil
}
