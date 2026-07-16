package main

import (
	"context"
	"fmt"
)

func (useCase *UserStatsUseCase) Get(ctx context.Context, userID int) (UserStats, error) {
	exec := useCase.db.Executor()

	if _, err := useCase.users.FindByID(ctx, exec, userID); err != nil {
		return UserStats{}, err
	}

	activeServices, err := useCase.services.CountActiveByProviderID(ctx, exec, userID)
	if err != nil {
		return UserStats{}, fmt.Errorf("fetch active services count: %w", err)
	}

	completedExchanges, err := useCase.exchanges.CountCompletedByUserID(ctx, exec, userID)
	if err != nil {
		return UserStats{}, fmt.Errorf("fetch completed exchanges count: %w", err)
	}

	creditBalance, err := useCase.creditTransactions.BalanceByUserID(ctx, exec, userID)
	if err != nil {
		return UserStats{}, fmt.Errorf("fetch credit balance: %w", err)
	}

	averageRating, reviewCount, err := useCase.reviews.StatsByTargetUserID(ctx, exec, userID)
	if err != nil {
		return UserStats{}, fmt.Errorf("fetch review stats: %w", err)
	}

	totalEarned, err := useCase.creditTransactions.TotalEarnedByUserID(ctx, exec, userID)
	if err != nil {
		return UserStats{}, fmt.Errorf("fetch total earned: %w", err)
	}

	totalSpent, err := useCase.creditTransactions.TotalSpentByUserID(ctx, exec, userID)
	if err != nil {
		return UserStats{}, fmt.Errorf("fetch total spent: %w", err)
	}

	return UserStats{
		UserID:             userID,
		ActiveServices:     activeServices,
		CompletedExchanges: completedExchanges,
		CreditBalance:      creditBalance,
		AverageRating:      averageRating,
		ReviewCount:        reviewCount,
		TotalEarned:        totalEarned,
		TotalSpent:         totalSpent,
	}, nil
}
