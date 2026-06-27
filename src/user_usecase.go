package main

import (
	"context"
	"fmt"
)

const welcomeCredits = 10

type userRepository interface {
	Create(ctx context.Context, exec dbExecutor, user User) (User, error)
	FindByID(ctx context.Context, exec dbExecutor, id int) (User, error)
	Update(ctx context.Context, exec dbExecutor, user User) (User, error)
}

type creditTransactionRepository interface {
	Create(ctx context.Context, exec dbExecutor, transaction CreditTransaction) error
	BalanceByUserID(ctx context.Context, exec dbExecutor, userID int) (int, error)
}

type database interface {
	Executor() dbExecutor
	WithinTransaction(ctx context.Context, fn func(exec dbExecutor) error) error
}

type UserUseCase struct {
	db                 database
	users              userRepository
	creditTransactions creditTransactionRepository
}

func NewUserUseCase(db database, users userRepository, creditTransactions creditTransactionRepository) *UserUseCase {
	return &UserUseCase{
		db:                 db,
		users:              users,
		creditTransactions: creditTransactions,
	}
}

func (useCase *UserUseCase) Register(ctx context.Context, pseudo, bio, ville string) (User, error) {
	user, err := NewUser(pseudo, bio, ville)
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
			Montant:   welcomeCredits,
			Type:      "earn",
			CreatedAt: user.CreatedAt,
		}
		return useCase.creditTransactions.Create(ctx, exec, welcomeTransaction)
	})
	if err != nil {
		return User{}, fmt.Errorf("creation de l'utilisateur : %w", err)
	}

	return user, nil
}

func (useCase *UserUseCase) Authenticate(ctx context.Context, id int) (User, error) {
	exec := useCase.db.Executor()

	return useCase.users.FindByID(ctx, exec, id)
}

func (useCase *UserUseCase) UpdateProfile(ctx context.Context, actorID, targetID int, pseudo, bio, ville string) (User, error) {
	if actorID != targetID {
		return User{}, ErrForbidden
	}

	changes, err := NewUser(pseudo, bio, ville)
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
	user.Ville = changes.Ville

	user, err = useCase.users.Update(ctx, exec, user)
	if err != nil {
		return User{}, fmt.Errorf("mise a jour du profil : %w", err)
	}

	balance, err := useCase.creditTransactions.BalanceByUserID(ctx, exec, targetID)
	if err != nil {
		return User{}, fmt.Errorf("calcul du solde : %w", err)
	}
	user.CreditBalance = balance

	return user, nil
}

func (useCase *UserUseCase) GetProfile(ctx context.Context, id int) (User, error) {
	exec := useCase.db.Executor()

	user, err := useCase.users.FindByID(ctx, exec, id)
	if err != nil {
		return User{}, err
	}

	balance, err := useCase.creditTransactions.BalanceByUserID(ctx, exec, id)
	if err != nil {
		return User{}, fmt.Errorf("calcul du solde : %w", err)
	}
	user.CreditBalance = balance

	return user, nil
}
