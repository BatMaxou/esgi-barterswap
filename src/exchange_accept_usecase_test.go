package main

import (
	"context"
	"errors"
	"testing"
)

func TestExchangeUseCaseAccept(t *testing.T) {
	t.Run("owner accepts pending exchange and blocks credits", func(t *testing.T) {
		exchanges := &fakeExchangeRepository{exchange: Exchange{
			ID: 5, ServiceID: 3, RequesterID: 2, OwnerID: 1, Status: ExchangeStatusPending,
		}}
		services := &fakeServiceRepository{service: Service{ID: 3, Credits: 2}}
		creditTransactions := &fakeCreditTransactionRepository{balance: 10}
		useCase := NewExchangeUseCase(&fakeDatabase{}, exchanges, services, creditTransactions)

		exchange, err := useCase.Accept(context.Background(), 1, 5)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if exchange.Status != ExchangeStatusAccepted {
			t.Errorf("Status = %q, want accepted", exchange.Status)
		}
		if !creditTransactions.createCalled {
			t.Fatal("credit transaction must be created")
		}
		if creditTransactions.transaction.Amount != -2 {
			t.Errorf("transaction amount = %d, want -2", creditTransactions.transaction.Amount)
		}
		if creditTransactions.transaction.Type != "spend" {
			t.Errorf("transaction type = %q, want spend", creditTransactions.transaction.Type)
		}
	})

	t.Run("non-owner -> ErrForbidden", func(t *testing.T) {
		exchanges := &fakeExchangeRepository{exchange: Exchange{ID: 5, OwnerID: 1, Status: ExchangeStatusPending}}
		useCase := NewExchangeUseCase(&fakeDatabase{}, exchanges, &fakeServiceRepository{}, &fakeCreditTransactionRepository{})

		_, err := useCase.Accept(context.Background(), 2, 5)
		if !errors.Is(err, ErrForbidden) {
			t.Fatalf("error = %v, want ErrForbidden", err)
		}
	})
}
