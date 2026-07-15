package main

import "context"

func (useCase *ExchangeUseCase) List(ctx context.Context, actorID int, status string) ([]Exchange, error) {
	if status != "" && !IsValidExchangeStatus(status) {
		return nil, ErrExchangeStatusInvalid
	}

	filter := ExchangeFilter{
		UserID: actorID,
		Status: status,
	}

	return useCase.exchanges.List(ctx, useCase.db.Executor(), filter)
}
