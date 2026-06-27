package main

import (
	"context"
	"errors"
	"testing"
)

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

func TestUserUseCaseAuthenticate(t *testing.T) {
	t.Run("utilisateur existant", func(t *testing.T) {
		users := &fakeUserRepository{user: User{ID: 5, Pseudo: "Thierry"}}
		useCase := NewUserUseCase(&fakeDatabase{}, users, &fakeCreditTransactionRepository{})

		user, err := useCase.Authenticate(context.Background(), 5)
		if err != nil {
			t.Fatalf("erreur inattendue : %v", err)
		}
		if user.ID != 5 {
			t.Errorf("ID = %d, attendu 5", user.ID)
		}
	})

	t.Run("utilisateur introuvable", func(t *testing.T) {
		users := &fakeUserRepository{findErr: ErrUserNotFound}
		useCase := NewUserUseCase(&fakeDatabase{}, users, &fakeCreditTransactionRepository{})

		_, err := useCase.Authenticate(context.Background(), 999)
		if !errors.Is(err, ErrUserNotFound) {
			t.Fatalf("erreur = %v, attendue ErrUserNotFound", err)
		}
	})
}

func TestUserUseCaseUpdateProfile(t *testing.T) {
	t.Run("mise a jour de son propre profil", func(t *testing.T) {
		users := &fakeUserRepository{user: User{ID: 5, Pseudo: "Ancien", CreatedAt: "2026-01-01T00:00:00Z"}}
		creditTransactions := &fakeCreditTransactionRepository{balance: 12}
		useCase := NewUserUseCase(&fakeDatabase{}, users, creditTransactions)

		user, err := useCase.UpdateProfile(context.Background(), 5, 5, "  Thierry  ", "nouvelle bio", "Lyon")
		if err != nil {
			t.Fatalf("erreur inattendue : %v", err)
		}
		if !users.updateCalled {
			t.Error("le repository Update doit etre appele")
		}
		if user.Pseudo != "Thierry" {
			t.Errorf("Pseudo = %q, attendu Thierry (trim applique)", user.Pseudo)
		}
		if user.Ville != "Lyon" {
			t.Errorf("Ville = %q, attendu Lyon", user.Ville)
		}
		if user.CreditBalance != 12 {
			t.Errorf("CreditBalance = %d, attendu 12 (solde recalcule)", user.CreditBalance)
		}
		if user.CreatedAt != "2026-01-01T00:00:00Z" {
			t.Errorf("CreatedAt = %q, doit etre preserve", user.CreatedAt)
		}
	})

	t.Run("modifier le profil d'un autre utilisateur renvoie ErrForbidden", func(t *testing.T) {
		users := &fakeUserRepository{user: User{ID: 5}}
		useCase := NewUserUseCase(&fakeDatabase{}, users, &fakeCreditTransactionRepository{})

		_, err := useCase.UpdateProfile(context.Background(), 5, 9, "Thierry", "", "")
		if !errors.Is(err, ErrForbidden) {
			t.Fatalf("erreur = %v, attendue ErrForbidden", err)
		}
		if users.updateCalled {
			t.Error("aucune ecriture ne doit avoir lieu en cas d'acces interdit")
		}
	})

	t.Run("pseudo vide renvoie ErrPseudoRequired", func(t *testing.T) {
		users := &fakeUserRepository{user: User{ID: 5}}
		useCase := NewUserUseCase(&fakeDatabase{}, users, &fakeCreditTransactionRepository{})

		_, err := useCase.UpdateProfile(context.Background(), 5, 5, "   ", "", "")
		if !errors.Is(err, ErrPseudoRequired) {
			t.Fatalf("erreur = %v, attendue ErrPseudoRequired", err)
		}
		if users.updateCalled {
			t.Error("aucune ecriture ne doit avoir lieu quand la validation echoue")
		}
	})
}
