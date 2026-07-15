package main

import (
	"context"
	"errors"
	"testing"
)

func TestUserUseCaseUpdateProfile(t *testing.T) {
	t.Run("updating your own profile", func(t *testing.T) {
		users := &fakeUserRepository{user: User{ID: 5, Pseudo: "Ancien", CreatedAt: "2026-01-01T00:00:00Z"}}
		creditTransactions := &fakeCreditTransactionRepository{balance: 12}
		useCase := NewUserUseCase(&fakeDatabase{}, users, creditTransactions)

		user, err := useCase.UpdateProfile(context.Background(), 5, 5, "  Thierry  ", "new bio", "Lyon")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !users.updateCalled {
			t.Error("the Update repository must be called")
		}
		if user.Pseudo != "Thierry" {
			t.Errorf("Pseudo = %q, want Thierry (trim applied)", user.Pseudo)
		}
		if user.City != "Lyon" {
			t.Errorf("City = %q, want Lyon", user.City)
		}
		if user.CreditBalance != 12 {
			t.Errorf("CreditBalance = %d, want 12 (recomputed balance)", user.CreditBalance)
		}
		if user.CreatedAt != "2026-01-01T00:00:00Z" {
			t.Errorf("CreatedAt = %q, must be preserved", user.CreatedAt)
		}
	})

	t.Run("updating another user's profile returns ErrForbidden", func(t *testing.T) {
		users := &fakeUserRepository{user: User{ID: 5}}
		useCase := NewUserUseCase(&fakeDatabase{}, users, &fakeCreditTransactionRepository{})

		_, err := useCase.UpdateProfile(context.Background(), 5, 9, "Thierry", "", "")
		if !errors.Is(err, ErrForbidden) {
			t.Fatalf("error = %v, want ErrForbidden", err)
		}
		if users.updateCalled {
			t.Error("no write must happen on a forbidden access")
		}
	})

	t.Run("empty pseudo returns ErrPseudoRequired", func(t *testing.T) {
		users := &fakeUserRepository{user: User{ID: 5}}
		useCase := NewUserUseCase(&fakeDatabase{}, users, &fakeCreditTransactionRepository{})

		_, err := useCase.UpdateProfile(context.Background(), 5, 5, "   ", "", "")
		if !errors.Is(err, ErrPseudoRequired) {
			t.Fatalf("error = %v, want ErrPseudoRequired", err)
		}
		if users.updateCalled {
			t.Error("no write must happen when validation fails")
		}
	})
}
