package main

import (
	"context"
	"errors"
	"testing"
)

func TestUserUseCaseAuthenticate(t *testing.T) {
	t.Run("existing user", func(t *testing.T) {
		users := &fakeUserRepository{user: User{ID: 5, Pseudo: "Thierry"}}
		useCase := NewUserUseCase(&fakeDatabase{}, users, &fakeCreditTransactionRepository{})

		user, err := useCase.Authenticate(context.Background(), 5)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if user.ID != 5 {
			t.Errorf("ID = %d, want 5", user.ID)
		}
	})

	t.Run("user not found", func(t *testing.T) {
		users := &fakeUserRepository{findErr: ErrUserNotFound}
		useCase := NewUserUseCase(&fakeDatabase{}, users, &fakeCreditTransactionRepository{})

		_, err := useCase.Authenticate(context.Background(), 999)
		if !errors.Is(err, ErrUserNotFound) {
			t.Fatalf("error = %v, want ErrUserNotFound", err)
		}
	})
}
