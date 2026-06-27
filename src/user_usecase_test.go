package main

import (
	"context"
	"errors"
	"testing"
)

type fakeUserRepository struct {
	createCalled bool
	user         User
	createErr    error
	findErr      error
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

func TestUserUseCaseRegister(t *testing.T) {
	t.Run("inscription valide attribue les credits de bienvenue", func(t *testing.T) {
		users := &fakeUserRepository{}
		creditTransactions := &fakeCreditTransactionRepository{}
		useCase := NewUserUseCase(&fakeDatabase{}, users, creditTransactions)

		user, err := useCase.Register(context.Background(), "  Thierry  ", "bio", "Paris")
		if err != nil {
			t.Fatalf("erreur inattendue : %v", err)
		}
		if user.ID != 7 {
			t.Errorf("ID = %d, attendu 7", user.ID)
		}
		if user.Pseudo != "Thierry" {
			t.Errorf("Pseudo = %q, attendu Thierry (trim applique)", user.Pseudo)
		}
		if user.CreditBalance != welcomeCredits {
			t.Errorf("CreditBalance = %d, attendu %d", user.CreditBalance, welcomeCredits)
		}
		if creditTransactions.transaction.UserID != 7 {
			t.Errorf("transaction UserID = %d, attendu 7", creditTransactions.transaction.UserID)
		}
		if creditTransactions.transaction.Montant != welcomeCredits {
			t.Errorf("montant de la transaction = %d, attendu %d", creditTransactions.transaction.Montant, welcomeCredits)
		}
		if creditTransactions.transaction.Type != "earn" {
			t.Errorf("type de la transaction = %q, attendu \"earn\"", creditTransactions.transaction.Type)
		}
	})

	t.Run("pseudo vide renvoie ErrPseudoRequired sans toucher aux repositories", func(t *testing.T) {
		users := &fakeUserRepository{}
		creditTransactions := &fakeCreditTransactionRepository{}
		useCase := NewUserUseCase(&fakeDatabase{}, users, creditTransactions)

		_, err := useCase.Register(context.Background(), "   ", "", "")
		if !errors.Is(err, ErrPseudoRequired) {
			t.Fatalf("erreur = %v, attendue ErrPseudoRequired", err)
		}
		if users.createCalled || creditTransactions.createCalled {
			t.Error("aucun repository ne doit etre appele quand la validation echoue")
		}
	})
}

func TestUserUseCaseGetProfile(t *testing.T) {
	t.Run("profil existant agrege le solde depuis le journal", func(t *testing.T) {
		users := &fakeUserRepository{user: User{ID: 5, Pseudo: "Thierry", CreatedAt: "2026-01-01T00:00:00Z"}}
		creditTransactions := &fakeCreditTransactionRepository{balance: 35}
		useCase := NewUserUseCase(&fakeDatabase{}, users, creditTransactions)

		user, err := useCase.GetProfile(context.Background(), 5)
		if err != nil {
			t.Fatalf("erreur inattendue : %v", err)
		}
		if user.ID != 5 {
			t.Errorf("ID = %d, attendu 5", user.ID)
		}
		if user.CreditBalance != 35 {
			t.Errorf("CreditBalance = %d, attendu 35 (somme du journal)", user.CreditBalance)
		}
	})

	t.Run("utilisateur introuvable", func(t *testing.T) {
		users := &fakeUserRepository{findErr: ErrUserNotFound}
		creditTransactions := &fakeCreditTransactionRepository{}
		useCase := NewUserUseCase(&fakeDatabase{}, users, creditTransactions)

		_, err := useCase.GetProfile(context.Background(), 999)
		if !errors.Is(err, ErrUserNotFound) {
			t.Fatalf("erreur = %v, attendue ErrUserNotFound", err)
		}
	})
}
