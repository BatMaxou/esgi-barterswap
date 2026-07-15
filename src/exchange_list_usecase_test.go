package main

import (
	"context"
	"errors"
	"testing"
)

func TestExchangeUseCaseList(t *testing.T) {
	t.Run("forwards actor and status filter", func(t *testing.T) {
		exchanges := &fakeExchangeRepository{exchanges: []Exchange{{ID: 1, Status: ExchangeStatusPending}}}
		useCase := NewExchangeUseCase(&fakeDatabase{}, exchanges, &fakeServiceRepository{}, &fakeCreditTransactionRepository{})

		got, err := useCase.List(context.Background(), 2, ExchangeStatusPending)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got) != 1 {
			t.Fatalf("exchanges = %d, want 1", len(got))
		}
		if exchanges.filter.UserID != 2 || exchanges.filter.Status != ExchangeStatusPending {
			t.Errorf("filter = %+v, want user 2 and pending status", exchanges.filter)
		}
	})

	t.Run("invalid status -> ErrExchangeStatusInvalid", func(t *testing.T) {
		useCase := NewExchangeUseCase(&fakeDatabase{}, &fakeExchangeRepository{}, &fakeServiceRepository{}, &fakeCreditTransactionRepository{})

		_, err := useCase.List(context.Background(), 2, "unknown")
		if !errors.Is(err, ErrExchangeStatusInvalid) {
			t.Fatalf("error = %v, want ErrExchangeStatusInvalid", err)
		}
	})
}
