package main

import (
	"context"
	"errors"
	"testing"
)

func TestExchangeUseCaseCreate(t *testing.T) {
	t.Run("valid request creates pending exchange", func(t *testing.T) {
		exchanges := &fakeExchangeRepository{}
		services := &fakeServiceRepository{service: Service{ID: 3, ProviderID: 1, Active: true, Credits: 2}}
		creditTransactions := &fakeCreditTransactionRepository{balance: 10}
		useCase := NewExchangeUseCase(&fakeDatabase{}, exchanges, services, creditTransactions)

		exchange, err := useCase.Create(context.Background(), 2, 3)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if exchange.Status != ExchangeStatusPending {
			t.Errorf("Status = %q, want pending", exchange.Status)
		}
		if !exchanges.createCalled {
			t.Error("exchange repository Create must be called")
		}
		if !exchanges.hasActiveCalled {
			t.Error("HasActiveForService must be called")
		}
	})

	t.Run("self request -> ErrExchangeSelfRequest without write", func(t *testing.T) {
		exchanges := &fakeExchangeRepository{}
		services := &fakeServiceRepository{service: Service{ID: 3, ProviderID: 2, Active: true, Credits: 2}}
		useCase := NewExchangeUseCase(&fakeDatabase{}, exchanges, services, &fakeCreditTransactionRepository{})

		_, err := useCase.Create(context.Background(), 2, 3)
		if !errors.Is(err, ErrExchangeSelfRequest) {
			t.Fatalf("error = %v, want ErrExchangeSelfRequest", err)
		}
		if exchanges.createCalled {
			t.Error("no repository write must happen")
		}
	})

	t.Run("insufficient credits -> ErrExchangeInsufficientCredits", func(t *testing.T) {
		exchanges := &fakeExchangeRepository{}
		services := &fakeServiceRepository{service: Service{ID: 3, ProviderID: 1, Active: true, Credits: 20}}
		creditTransactions := &fakeCreditTransactionRepository{balance: 10}
		useCase := NewExchangeUseCase(&fakeDatabase{}, exchanges, services, creditTransactions)

		_, err := useCase.Create(context.Background(), 2, 3)
		if !errors.Is(err, ErrExchangeInsufficientCredits) {
			t.Fatalf("error = %v, want ErrExchangeInsufficientCredits", err)
		}
		if exchanges.createCalled {
			t.Error("exchange must not be created")
		}
	})

	t.Run("active exchange on service -> ErrExchangeServiceUnavailable", func(t *testing.T) {
		exchanges := &fakeExchangeRepository{hasActive: true}
		services := &fakeServiceRepository{service: Service{ID: 3, ProviderID: 1, Active: true, Credits: 2}}
		useCase := NewExchangeUseCase(&fakeDatabase{}, exchanges, services, &fakeCreditTransactionRepository{balance: 10})

		_, err := useCase.Create(context.Background(), 2, 3)
		if !errors.Is(err, ErrExchangeServiceUnavailable) {
			t.Fatalf("error = %v, want ErrExchangeServiceUnavailable", err)
		}
	})

	t.Run("inactive service -> ErrExchangeServiceInactive", func(t *testing.T) {
		services := &fakeServiceRepository{service: Service{ID: 3, ProviderID: 1, Active: false, Credits: 2}}
		useCase := NewExchangeUseCase(&fakeDatabase{}, &fakeExchangeRepository{}, services, &fakeCreditTransactionRepository{})

		_, err := useCase.Create(context.Background(), 2, 3)
		if !errors.Is(err, ErrExchangeServiceInactive) {
			t.Fatalf("error = %v, want ErrExchangeServiceInactive", err)
		}
	})
}
