package main

import (
	"context"
	"errors"
	"testing"
)

func TestReviewUseCaseListForUser(t *testing.T) {
	t.Run("existing user -> reviews", func(t *testing.T) {
		reviews := &fakeReviewRepository{userReviews: []Review{{ID: 1, TargetID: 1, Rating: 5}}}
		useCase := NewReviewUseCase(&fakeDatabase{}, reviews, &fakeExchangeRepository{}, &fakeUserRepository{user: User{ID: 1}}, &fakeServiceRepository{})

		got, err := useCase.ListForUser(context.Background(), 1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got) != 1 {
			t.Fatalf("reviews = %d, want 1", len(got))
		}
	})

	t.Run("unknown user -> ErrUserNotFound", func(t *testing.T) {
		useCase := NewReviewUseCase(&fakeDatabase{}, &fakeReviewRepository{}, &fakeExchangeRepository{}, &fakeUserRepository{findErr: ErrUserNotFound}, &fakeServiceRepository{})

		_, err := useCase.ListForUser(context.Background(), 999)
		if !errors.Is(err, ErrUserNotFound) {
			t.Fatalf("error = %v, want ErrUserNotFound", err)
		}
	})
}
