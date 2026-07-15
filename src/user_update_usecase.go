package main

import (
	"context"
	"fmt"
)

func (useCase *UserUseCase) UpdateProfile(ctx context.Context, actorID, targetID int, pseudo, bio, city string) (User, error) {
	if actorID != targetID {
		return User{}, ErrForbidden
	}

	changes, err := NewUser(pseudo, bio, city)
	if err != nil {
		return User{}, err
	}

	exec := useCase.db.Executor()

	user, err := useCase.users.FindByID(ctx, exec, targetID)
	if err != nil {
		return User{}, err
	}

	user.Pseudo = changes.Pseudo
	user.Bio = changes.Bio
	user.City = changes.City

	user, err = useCase.users.Update(ctx, exec, user)
	if err != nil {
		return User{}, fmt.Errorf("update profile: %w", err)
	}

	balance, err := useCase.creditTransactions.BalanceByUserID(ctx, exec, targetID)
	if err != nil {
		return User{}, fmt.Errorf("compute balance: %w", err)
	}
	user.CreditBalance = balance

	return user, nil
}
