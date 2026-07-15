package main

import "context"

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
