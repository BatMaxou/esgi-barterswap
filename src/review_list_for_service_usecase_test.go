package main

import (
	"context"
	"errors"
	"testing"
)

func TestReviewUseCaseListForService(t *testing.T) {
	t.Run("existing service -> reviews from its exchanges", func(t *testing.T) {
		exchanges := &fakeExchangeRepository{serviceExchanges: []Exchange{{ID: 3, ServiceID: 7}}}
		reviews := &fakeReviewRepository{exchangeReviews: []Review{{ID: 1, ExchangeID: 3, Rating: 4}}}
		useCase := NewReviewUseCase(&fakeDatabase{}, reviews, exchanges, &fakeUserRepository{}, &fakeServiceRepository{service: Service{ID: 7}})

		got, err := useCase.ListForService(context.Background(), 7)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got) != 1 {
			t.Fatalf("reviews = %d, want 1", len(got))
		}
	})

	t.Run("unknown service -> ErrServiceNotFound", func(t *testing.T) {
		useCase := NewReviewUseCase(&fakeDatabase{}, &fakeReviewRepository{}, &fakeExchangeRepository{}, &fakeUserRepository{}, &fakeServiceRepository{findErr: ErrServiceNotFound})

		_, err := useCase.ListForService(context.Background(), 999)
		if !errors.Is(err, ErrServiceNotFound) {
			t.Fatalf("error = %v, want ErrServiceNotFound", err)
		}
	})

	t.Run("service without exchanges -> empty list", func(t *testing.T) {
		useCase := NewReviewUseCase(&fakeDatabase{}, &fakeReviewRepository{}, &fakeExchangeRepository{}, &fakeUserRepository{}, &fakeServiceRepository{service: Service{ID: 7}})

		got, err := useCase.ListForService(context.Background(), 7)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got) != 0 {
			t.Fatalf("reviews = %d, want 0", len(got))
		}
	})
}
