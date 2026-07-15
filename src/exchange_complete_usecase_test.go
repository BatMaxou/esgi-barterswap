package main

import (
	"context"
	"errors"
	"testing"
)

func TestExchangeUseCaseComplete(t *testing.T) {
	t.Run("owner completes accepted exchange and credits provider", func(t *testing.T) {
		exchanges := &fakeExchangeRepository{exchange: Exchange{
			ID: 5, ServiceID: 3, OwnerID: 1, Status: ExchangeStatusAccepted,
		}}
		services := &fakeServiceRepository{service: Service{ID: 3, Credits: 2}}
		creditTransactions := &fakeCreditTransactionRepository{}
		useCase := NewExchangeUseCase(&fakeDatabase{}, exchanges, services, creditTransactions)

		exchange, err := useCase.Complete(context.Background(), 1, 5)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if exchange.Status != ExchangeStatusCompleted {
			t.Errorf("Status = %q, want completed", exchange.Status)
		}
		if creditTransactions.transaction.UserID != 1 {
			t.Errorf("transaction UserID = %d, want 1", creditTransactions.transaction.UserID)
		}
		if creditTransactions.transaction.Amount != 2 {
			t.Errorf("transaction amount = %d, want 2", creditTransactions.transaction.Amount)
		}
		if creditTransactions.transaction.Type != "earn" {
			t.Errorf("transaction type = %q, want earn", creditTransactions.transaction.Type)
		}
	})

	t.Run("pending exchange -> ErrExchangeInvalidTransition", func(t *testing.T) {
		exchanges := &fakeExchangeRepository{exchange: Exchange{ID: 5, OwnerID: 1, Status: ExchangeStatusPending}}
		useCase := NewExchangeUseCase(&fakeDatabase{}, exchanges, &fakeServiceRepository{}, &fakeCreditTransactionRepository{})

		_, err := useCase.Complete(context.Background(), 1, 5)
		if !errors.Is(err, ErrExchangeInvalidTransition) {
			t.Fatalf("error = %v, want ErrExchangeInvalidTransition", err)
		}
	})
}
