package main

import "context"

// Fakes shared by the UserUseCase operation tests (and reused by the
// SkillUseCase tests).

type fakeUserRepository struct {
	createCalled bool
	updateCalled bool
	user         User
	updatedUser  User
	createErr    error
	findErr      error
	updateErr    error
}

func (fake *fakeUserRepository) Create(ctx context.Context, exec dbExecutor, user User) (User, error) {
	fake.createCalled = true
	if fake.createErr != nil {
		return User{}, fake.createErr
	}
	user.ID = 7
	return user, nil
}

func (fake *fakeUserRepository) FindByID(ctx context.Context, exec dbExecutor, id int) (User, error) {
	if fake.findErr != nil {
		return User{}, fake.findErr
	}
	return fake.user, nil
}

func (fake *fakeUserRepository) Update(ctx context.Context, exec dbExecutor, user User) (User, error) {
	fake.updateCalled = true
	if fake.updateErr != nil {
		return User{}, fake.updateErr
	}
	fake.updatedUser = user
	return user, nil
}

type fakeCreditTransactionRepository struct {
	createCalled bool
	transaction  CreditTransaction
	balance      int
	createErr    error
}

func (fake *fakeCreditTransactionRepository) Create(ctx context.Context, exec dbExecutor, transaction CreditTransaction) error {
	fake.createCalled = true
	fake.transaction = transaction
	return fake.createErr
}

func (fake *fakeCreditTransactionRepository) BalanceByUserID(ctx context.Context, exec dbExecutor, userID int) (int, error) {
	return fake.balance, nil
}

type fakeDatabase struct{}

func (fake *fakeDatabase) Executor() dbExecutor { return nil }

func (fake *fakeDatabase) WithinTransaction(ctx context.Context, fn func(exec dbExecutor) error) error {
	return fn(nil)
}
