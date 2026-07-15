package main

import "context"

func (useCase *ExchangeUseCase) Get(ctx context.Context, actorID, exchangeID int) (Exchange, error) {
	exchange, err := useCase.exchanges.FindByID(ctx, useCase.db.Executor(), exchangeID)
	if err != nil {
		return Exchange{}, err
	}
	if err := useCase.ensureActorCanView(actorID, exchange); err != nil {
		return Exchange{}, err
	}

	return exchange, nil
}
