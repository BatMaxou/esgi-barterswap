package main

import (
	"context"
	"errors"
	"testing"
)

func TestUserUseCaseRegister(t *testing.T) {
	t.Run("valid registration grants welcome credits", func(t *testing.T) {
		users := &fakeUserRepository{}
		creditTransactions := &fakeCreditTransactionRepository{}
		useCase := NewUserUseCase(&fakeDatabase{}, users, creditTransactions)

		user, err := useCase.Register(context.Background(), "  Thierry  ", "bio", "Paris")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if user.ID != 7 {
			t.Errorf("ID = %d, want 7", user.ID)
		}
		if user.Pseudo != "Thierry" {
			t.Errorf("Pseudo = %q, want Thierry (trim applied)", user.Pseudo)
		}
		if user.CreditBalance != welcomeCredits {
			t.Errorf("CreditBalance = %d, want %d", user.CreditBalance, welcomeCredits)
		}
		if creditTransactions.transaction.UserID != 7 {
			t.Errorf("transaction UserID = %d, want 7", creditTransactions.transaction.UserID)
		}
		if creditTransactions.transaction.Amount != welcomeCredits {
			t.Errorf("transaction amount = %d, want %d", creditTransactions.transaction.Amount, welcomeCredits)
		}
		if creditTransactions.transaction.Type != "earn" {
			t.Errorf("transaction type = %q, want \"earn\"", creditTransactions.transaction.Type)
		}
	})

	t.Run("empty pseudo returns ErrPseudoRequired without touching repositories", func(t *testing.T) {
		users := &fakeUserRepository{}
		creditTransactions := &fakeCreditTransactionRepository{}
		useCase := NewUserUseCase(&fakeDatabase{}, users, creditTransactions)

		_, err := useCase.Register(context.Background(), "   ", "", "")
		if !errors.Is(err, ErrPseudoRequired) {
			t.Fatalf("error = %v, want ErrPseudoRequired", err)
		}
		if users.createCalled || creditTransactions.createCalled {
			t.Error("no repository must be called when validation fails")
		}
	})
}
