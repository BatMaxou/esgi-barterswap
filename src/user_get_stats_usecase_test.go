package main

import (
	"context"
	"errors"
	"testing"
)

func TestUserStatsUseCaseGet(t *testing.T) {
	t.Run("aggregates all metrics for an existing user", func(t *testing.T) {
		useCase := NewUserStatsUseCase(
			&fakeDatabase{},
			&fakeUserRepository{user: User{ID: 1}},
			&fakeServiceRepository{activeCount: 2},
			&fakeExchangeRepository{completedCount: 3},
			&fakeReviewRepository{averageRating: 4.5, reviewCount: 2},
			&fakeCreditTransactionRepository{balance: 8, totalEarned: 12, totalSpent: 2},
		)

		stats, err := useCase.Get(context.Background(), 1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if stats.UserID != 1 {
			t.Errorf("UserID = %d, want 1", stats.UserID)
		}
		if stats.ActiveServices != 2 {
			t.Errorf("ActiveServices = %d, want 2", stats.ActiveServices)
		}
		if stats.CompletedExchanges != 3 {
			t.Errorf("CompletedExchanges = %d, want 3", stats.CompletedExchanges)
		}
		if stats.CreditBalance != 8 {
			t.Errorf("CreditBalance = %d, want 8", stats.CreditBalance)
		}
		if stats.AverageRating != 4.5 {
			t.Errorf("AverageRating = %v, want 4.5", stats.AverageRating)
		}
		if stats.ReviewCount != 2 {
			t.Errorf("ReviewCount = %d, want 2", stats.ReviewCount)
		}
		if stats.TotalEarned != 12 {
			t.Errorf("TotalEarned = %d, want 12", stats.TotalEarned)
		}
		if stats.TotalSpent != 2 {
			t.Errorf("TotalSpent = %d, want 2", stats.TotalSpent)
		}
	})

	t.Run("unknown user -> ErrUserNotFound", func(t *testing.T) {
		useCase := NewUserStatsUseCase(
			&fakeDatabase{},
			&fakeUserRepository{findErr: ErrUserNotFound},
			&fakeServiceRepository{},
			&fakeExchangeRepository{},
			&fakeReviewRepository{},
			&fakeCreditTransactionRepository{},
		)

		_, err := useCase.Get(context.Background(), 999)
		if !errors.Is(err, ErrUserNotFound) {
			t.Fatalf("error = %v, want ErrUserNotFound", err)
		}
	})

	t.Run("user without reviews -> zero rating stats", func(t *testing.T) {
		useCase := NewUserStatsUseCase(
			&fakeDatabase{},
			&fakeUserRepository{user: User{ID: 1}},
			&fakeServiceRepository{},
			&fakeExchangeRepository{},
			&fakeReviewRepository{},
			&fakeCreditTransactionRepository{balance: 10, totalEarned: 10},
		)

		stats, err := useCase.Get(context.Background(), 1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if stats.AverageRating != 0 || stats.ReviewCount != 0 {
			t.Errorf("AverageRating = %v, ReviewCount = %d, want 0/0", stats.AverageRating, stats.ReviewCount)
		}
	})
}
