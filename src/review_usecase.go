package main

import "context"

type reviewRepository interface {
	Create(ctx context.Context, exec dbExecutor, review Review) (Review, error)
	ListByTargetUserID(ctx context.Context, exec dbExecutor, userID int) ([]Review, error)
	ListByExchangeIDs(ctx context.Context, exec dbExecutor, exchangeIDs []int) ([]Review, error)
}

type ReviewUseCase struct {
	db        database
	reviews   reviewRepository
	exchanges exchangeRepository
	users     userRepository
	services  serviceRepository
}

func NewReviewUseCase(
	db database,
	reviews reviewRepository,
	exchanges exchangeRepository,
	users userRepository,
	services serviceRepository,
) *ReviewUseCase {
	return &ReviewUseCase{
		db:        db,
		reviews:   reviews,
		exchanges: exchanges,
		users:     users,
		services:  services,
	}
}

func reviewTargetForAuthor(exchange Exchange, authorID int) (int, error) {
	switch authorID {
	case exchange.RequesterID:
		return exchange.OwnerID, nil
	case exchange.OwnerID:
		return exchange.RequesterID, nil
	default:
		return 0, ErrReviewNotParticipant
	}
}
