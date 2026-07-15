package main

import (
	"context"
	"errors"
	"testing"
)

func TestExchangeUseCaseCancel(t *testing.T) {
	t.Run("cancel pending exchange without refund", func(t *testing.T) {
		exchanges := &fakeExchangeRepository{exchange: Exchange{
			ID: 5, RequesterID: 2, OwnerID: 1, Status: ExchangeStatusPending,
		}}
		creditTransactions := &fakeCreditTransactionRepository{}
		useCase := NewExchangeUseCase(&fakeDatabase{}, exchanges, &fakeServiceRepository{}, creditTransactions)

		exchange, err := useCase.Cancel(context.Background(), 2, 5)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if exchange.Status != ExchangeStatusCancelled {
			t.Errorf("Status = %q, want cancelled", exchange.Status)
		}
		if creditTransactions.createCalled {
			t.Error("no credit transaction must be created for pending cancel")
		}
	})

	t.Run("cancel accepted exchange refunds requester", func(t *testing.T) {
		exchanges := &fakeExchangeRepository{exchange: Exchange{
			ID: 5, ServiceID: 3, RequesterID: 2, OwnerID: 1, Status: ExchangeStatusAccepted,
		}}
		services := &fakeServiceRepository{service: Service{ID: 3, Credits: 2}}
		creditTransactions := &fakeCreditTransactionRepository{}
		useCase := NewExchangeUseCase(&fakeDatabase{}, exchanges, services, creditTransactions)

		exchange, err := useCase.Cancel(context.Background(), 1, 5)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if exchange.Status != ExchangeStatusCancelled {
			t.Errorf("Status = %q, want cancelled", exchange.Status)
		}
		if creditTransactions.transaction.UserID != 2 {
			t.Errorf("transaction UserID = %d, want 2", creditTransactions.transaction.UserID)
		}
		if creditTransactions.transaction.Amount != 2 {
			t.Errorf("transaction amount = %d, want 2", creditTransactions.transaction.Amount)
		}
		if creditTransactions.transaction.Type != "refund" {
			t.Errorf("transaction type = %q, want refund", creditTransactions.transaction.Type)
		}
	})

	t.Run("completed exchange -> ErrExchangeInvalidTransition", func(t *testing.T) {
		exchanges := &fakeExchangeRepository{exchange: Exchange{ID: 5, RequesterID: 2, OwnerID: 1, Status: ExchangeStatusCompleted}}
		useCase := NewExchangeUseCase(&fakeDatabase{}, exchanges, &fakeServiceRepository{}, &fakeCreditTransactionRepository{})

		_, err := useCase.Cancel(context.Background(), 2, 5)
		if !errors.Is(err, ErrExchangeInvalidTransition) {
			t.Fatalf("error = %v, want ErrExchangeInvalidTransition", err)
		}
	})
}
