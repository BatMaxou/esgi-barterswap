package main

import "context"

func (useCase *ReviewUseCase) ListForService(ctx context.Context, serviceID int) ([]Review, error) {
	exec := useCase.db.Executor()

	if _, err := useCase.services.FindByID(ctx, exec, serviceID); err != nil {
		return nil, err
	}

	exchanges, err := useCase.exchanges.ListByServiceID(ctx, exec, serviceID)
	if err != nil {
		return nil, err
	}

	exchangeIDs := make([]int, 0, len(exchanges))
	for _, exchange := range exchanges {
		exchangeIDs = append(exchangeIDs, exchange.ID)
	}

	return useCase.reviews.ListByExchangeIDs(ctx, exec, exchangeIDs)
}
