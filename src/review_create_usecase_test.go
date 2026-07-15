package main

import (
	"context"
	"errors"
	"testing"
)

func TestReviewUseCaseCreate(t *testing.T) {
	t.Run("participant reviews completed exchange -> 201 path", func(t *testing.T) {
		exchanges := &fakeExchangeRepository{exchange: Exchange{
			ID: 5, RequesterID: 2, OwnerID: 1, Status: ExchangeStatusCompleted,
		}}
		reviews := &fakeReviewRepository{}
		useCase := NewReviewUseCase(&fakeDatabase{}, reviews, exchanges, &fakeUserRepository{}, &fakeServiceRepository{})

		review, err := useCase.Create(context.Background(), 2, 5, 5, "Great")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if review.TargetID != 1 {
			t.Errorf("TargetID = %d, want 1", review.TargetID)
		}
		if !reviews.createCalled {
			t.Error("review repository Create must be called")
		}
	})

	t.Run("pending exchange -> ErrReviewExchangeNotCompleted", func(t *testing.T) {
		exchanges := &fakeExchangeRepository{exchange: Exchange{
			ID: 5, RequesterID: 2, OwnerID: 1, Status: ExchangeStatusPending,
		}}
		useCase := NewReviewUseCase(&fakeDatabase{}, &fakeReviewRepository{}, exchanges, &fakeUserRepository{}, &fakeServiceRepository{})

		_, err := useCase.Create(context.Background(), 2, 5, 5, "Great")
		if !errors.Is(err, ErrReviewExchangeNotCompleted) {
			t.Fatalf("error = %v, want ErrReviewExchangeNotCompleted", err)
		}
	})

	t.Run("outsider -> ErrReviewNotParticipant", func(t *testing.T) {
		exchanges := &fakeExchangeRepository{exchange: Exchange{
			ID: 5, RequesterID: 2, OwnerID: 1, Status: ExchangeStatusCompleted,
		}}
		useCase := NewReviewUseCase(&fakeDatabase{}, &fakeReviewRepository{}, exchanges, &fakeUserRepository{}, &fakeServiceRepository{})

		_, err := useCase.Create(context.Background(), 99, 5, 5, "Great")
		if !errors.Is(err, ErrReviewNotParticipant) {
			t.Fatalf("error = %v, want ErrReviewNotParticipant", err)
		}
	})

	t.Run("duplicate review -> ErrReviewAlreadyExists", func(t *testing.T) {
		exchanges := &fakeExchangeRepository{exchange: Exchange{
			ID: 5, RequesterID: 2, OwnerID: 1, Status: ExchangeStatusCompleted,
		}}
		reviews := &fakeReviewRepository{createErr: ErrReviewAlreadyExists}
		useCase := NewReviewUseCase(&fakeDatabase{}, reviews, exchanges, &fakeUserRepository{}, &fakeServiceRepository{})

		_, err := useCase.Create(context.Background(), 2, 5, 5, "Great")
		if !errors.Is(err, ErrReviewAlreadyExists) {
			t.Fatalf("error = %v, want ErrReviewAlreadyExists", err)
		}
	})

	t.Run("invalid rating -> ErrReviewRatingInvalid without write", func(t *testing.T) {
		exchanges := &fakeExchangeRepository{exchange: Exchange{
			ID: 5, RequesterID: 2, OwnerID: 1, Status: ExchangeStatusCompleted,
		}}
		reviews := &fakeReviewRepository{}
		useCase := NewReviewUseCase(&fakeDatabase{}, reviews, exchanges, &fakeUserRepository{}, &fakeServiceRepository{})

		_, err := useCase.Create(context.Background(), 2, 5, 0, "Great")
		if !errors.Is(err, ErrReviewRatingInvalid) {
			t.Fatalf("error = %v, want ErrReviewRatingInvalid", err)
		}
		if reviews.createCalled {
			t.Error("no repository write must happen when validation fails")
		}
	})
}
