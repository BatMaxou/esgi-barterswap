package main

import (
	"context"
	"errors"
	"testing"
)

func TestExchangeUseCaseGet(t *testing.T) {
	t.Run("participant can view exchange", func(t *testing.T) {
		exchanges := &fakeExchangeRepository{exchange: Exchange{ID: 4, RequesterID: 2, OwnerID: 1}}
		useCase := NewExchangeUseCase(&fakeDatabase{}, exchanges, &fakeServiceRepository{}, &fakeCreditTransactionRepository{})

		exchange, err := useCase.Get(context.Background(), 2, 4)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if exchange.ID != 4 {
			t.Errorf("ID = %d, want 4", exchange.ID)
		}
	})

	t.Run("outsider -> ErrForbidden", func(t *testing.T) {
		exchanges := &fakeExchangeRepository{exchange: Exchange{ID: 4, RequesterID: 2, OwnerID: 1}}
		useCase := NewExchangeUseCase(&fakeDatabase{}, exchanges, &fakeServiceRepository{}, &fakeCreditTransactionRepository{})

		_, err := useCase.Get(context.Background(), 99, 4)
		if !errors.Is(err, ErrForbidden) {
			t.Fatalf("error = %v, want ErrForbidden", err)
		}
	})
}
