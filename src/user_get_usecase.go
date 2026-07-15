package main

import (
	"context"
	"fmt"
)

func (useCase *UserUseCase) GetProfile(ctx context.Context, id int) (User, error) {
	exec := useCase.db.Executor()

	user, err := useCase.users.FindByID(ctx, exec, id)
	if err != nil {
		return User{}, err
	}

	balance, err := useCase.creditTransactions.BalanceByUserID(ctx, exec, id)
	if err != nil {
		return User{}, fmt.Errorf("compute balance: %w", err)
	}
	user.CreditBalance = balance

	return user, nil
}
