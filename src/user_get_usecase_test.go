package main

import (
	"context"
	"errors"
	"testing"
)

func TestUserUseCaseGetProfile(t *testing.T) {
	t.Run("existing profile aggregates the balance from the ledger", func(t *testing.T) {
		users := &fakeUserRepository{user: User{ID: 5, Pseudo: "Thierry", CreatedAt: "2026-01-01T00:00:00Z"}}
		creditTransactions := &fakeCreditTransactionRepository{balance: 35}
		useCase := NewUserUseCase(&fakeDatabase{}, users, creditTransactions)

		user, err := useCase.GetProfile(context.Background(), 5)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if user.ID != 5 {
			t.Errorf("ID = %d, want 5", user.ID)
		}
		if user.CreditBalance != 35 {
			t.Errorf("CreditBalance = %d, want 35 (ledger sum)", user.CreditBalance)
		}
	})

	t.Run("user not found", func(t *testing.T) {
		users := &fakeUserRepository{findErr: ErrUserNotFound}
		creditTransactions := &fakeCreditTransactionRepository{}
		useCase := NewUserUseCase(&fakeDatabase{}, users, creditTransactions)

		_, err := useCase.GetProfile(context.Background(), 999)
		if !errors.Is(err, ErrUserNotFound) {
			t.Fatalf("error = %v, want ErrUserNotFound", err)
		}
	})
}
