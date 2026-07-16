package main

import "context"

type userStatsUserRepository interface {
	FindByID(ctx context.Context, exec dbExecutor, id int) (User, error)
}

type userStatsServiceRepository interface {
	CountActiveByProviderID(ctx context.Context, exec dbExecutor, providerID int) (int, error)
}

type userStatsExchangeRepository interface {
	CountCompletedByUserID(ctx context.Context, exec dbExecutor, userID int) (int, error)
}

type userStatsReviewRepository interface {
	StatsByTargetUserID(ctx context.Context, exec dbExecutor, userID int) (float64, int, error)
}

type userStatsCreditTransactionRepository interface {
	BalanceByUserID(ctx context.Context, exec dbExecutor, userID int) (int, error)
	TotalEarnedByUserID(ctx context.Context, exec dbExecutor, userID int) (int, error)
	TotalSpentByUserID(ctx context.Context, exec dbExecutor, userID int) (int, error)
}

type UserStatsUseCase struct {
	db                 database
	users              userStatsUserRepository
	services           userStatsServiceRepository
	exchanges          userStatsExchangeRepository
	reviews            userStatsReviewRepository
	creditTransactions userStatsCreditTransactionRepository
}

func NewUserStatsUseCase(
	db database,
	users userStatsUserRepository,
	services userStatsServiceRepository,
	exchanges userStatsExchangeRepository,
	reviews userStatsReviewRepository,
	creditTransactions userStatsCreditTransactionRepository,
) *UserStatsUseCase {
	return &UserStatsUseCase{
		db:                 db,
		users:              users,
		services:           services,
		exchanges:          exchanges,
		reviews:            reviews,
		creditTransactions: creditTransactions,
	}
}
