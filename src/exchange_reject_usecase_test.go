package main

import (
	"context"
	"errors"
	"testing"
)

func TestExchangeUseCaseReject(t *testing.T) {
	t.Run("owner rejects pending exchange", func(t *testing.T) {
		exchanges := &fakeExchangeRepository{exchange: Exchange{ID: 5, OwnerID: 1, Status: ExchangeStatusPending}}
		useCase := NewExchangeUseCase(&fakeDatabase{}, exchanges, &fakeServiceRepository{}, &fakeCreditTransactionRepository{})

		exchange, err := useCase.Reject(context.Background(), 1, 5)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if exchange.Status != ExchangeStatusRejected {
			t.Errorf("Status = %q, want rejected", exchange.Status)
		}
	})

	t.Run("accepted exchange -> ErrExchangeInvalidTransition", func(t *testing.T) {
		exchanges := &fakeExchangeRepository{exchange: Exchange{ID: 5, OwnerID: 1, Status: ExchangeStatusAccepted}}
		useCase := NewExchangeUseCase(&fakeDatabase{}, exchanges, &fakeServiceRepository{}, &fakeCreditTransactionRepository{})

		_, err := useCase.Reject(context.Background(), 1, 5)
		if !errors.Is(err, ErrExchangeInvalidTransition) {
			t.Fatalf("error = %v, want ErrExchangeInvalidTransition", err)
		}
	})
}
