package main

import (
	"context"
	"fmt"
)

func (useCase *ReviewUseCase) Create(ctx context.Context, actorID, exchangeID, rating int, comment string) (Review, error) {
	exec := useCase.db.Executor()

	exchange, err := useCase.exchanges.FindByID(ctx, exec, exchangeID)
	if err != nil {
		return Review{}, err
	}
	if exchange.Status != ExchangeStatusCompleted {
		return Review{}, ErrReviewExchangeNotCompleted
	}

	targetID, err := reviewTargetForAuthor(exchange, actorID)
	if err != nil {
		return Review{}, err
	}

	review, err := NewReview(exchangeID, actorID, targetID, rating, comment)
	if err != nil {
		return Review{}, err
	}

	created, err := useCase.reviews.Create(ctx, exec, review)
	if err != nil {
		return Review{}, fmt.Errorf("create review: %w", err)
	}

	return created, nil
}
