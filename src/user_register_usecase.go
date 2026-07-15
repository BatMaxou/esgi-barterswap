package main

import (
	"context"
	"fmt"
)

func (useCase *UserUseCase) Register(ctx context.Context, pseudo, bio, city string) (User, error) {
	user, err := NewUser(pseudo, bio, city)
	if err != nil {
		return User{}, err
	}
	user.CreditBalance = welcomeCredits

	err = useCase.db.WithinTransaction(ctx, func(exec dbExecutor) error {
		created, err := useCase.users.Create(ctx, exec, user)
		if err != nil {
			return err
		}
		user = created

		welcomeTransaction := CreditTransaction{
			UserID:    user.ID,
			Amount:    welcomeCredits,
			Type:      "earn",
			CreatedAt: user.CreatedAt,
		}
		return useCase.creditTransactions.Create(ctx, exec, welcomeTransaction)
	})
	if err != nil {
		return User{}, fmt.Errorf("create user: %w", err)
	}

	return user, nil
}
