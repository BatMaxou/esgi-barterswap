package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type fakeReviewUseCase struct {
	createFunc         func(ctx context.Context, actorID, exchangeID, rating int, comment string) (Review, error)
	listForUserFunc    func(ctx context.Context, userID int) ([]Review, error)
	listForServiceFunc func(ctx context.Context, serviceID int) ([]Review, error)
}

func (fake *fakeReviewUseCase) Create(ctx context.Context, actorID, exchangeID, rating int, comment string) (Review, error) {
	return fake.createFunc(ctx, actorID, exchangeID, rating, comment)
}

func (fake *fakeReviewUseCase) ListForUser(ctx context.Context, userID int) ([]Review, error) {
	return fake.listForUserFunc(ctx, userID)
}

func (fake *fakeReviewUseCase) ListForService(ctx context.Context, serviceID int) ([]Review, error) {
	return fake.listForServiceFunc(ctx, serviceID)
}

func TestHandleCreateReview(t *testing.T) {
	t.Run("valid review -> 201", func(t *testing.T) {
		app := &api{reviews: &fakeReviewUseCase{
			createFunc: func(ctx context.Context, actorID, exchangeID, rating int, comment string) (Review, error) {
				return Review{ID: 1, ExchangeID: exchangeID, AuthorID: actorID, Rating: rating}, nil
			},
		}}

		req := httptest.NewRequest(http.MethodPost, "/api/exchanges/5/review", strings.NewReader(`{"rating":5,"comment":"Great"}`))
		req.SetPathValue("id", "5")
		req.Header.Set("Content-Type", "application/json")
		req = withUser(req, User{ID: 2})
		rec := httptest.NewRecorder()

		app.handleCreateReview(rec, req)

		if rec.Code != http.StatusCreated {
			t.Fatalf("status = %d, want 201", rec.Code)
		}
	})

	t.Run("exchange not completed -> 400", func(t *testing.T) {
		app := &api{reviews: &fakeReviewUseCase{
			createFunc: func(ctx context.Context, actorID, exchangeID, rating int, comment string) (Review, error) {
				return Review{}, ErrReviewExchangeNotCompleted
			},
		}}

		req := httptest.NewRequest(http.MethodPost, "/api/exchanges/5/review", strings.NewReader(`{"rating":5}`))
		req.SetPathValue("id", "5")
		req.Header.Set("Content-Type", "application/json")
		req = withUser(req, User{ID: 2})
		rec := httptest.NewRecorder()

		app.handleCreateReview(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want 400", rec.Code)
		}
	})
}

func TestHandleListUserReviews(t *testing.T) {
	t.Run("existing user -> 200", func(t *testing.T) {
		app := &api{reviews: &fakeReviewUseCase{
			listForUserFunc: func(ctx context.Context, userID int) ([]Review, error) {
				return []Review{{ID: 1, TargetID: userID, Rating: 5}}, nil
			},
		}}

		req := httptest.NewRequest(http.MethodGet, "/api/users/1/reviews", nil)
		req.SetPathValue("id", "1")
		rec := httptest.NewRecorder()

		app.handleListUserReviews(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200", rec.Code)
		}
	})
}

func TestHandleListServiceReviews(t *testing.T) {
	t.Run("existing service -> 200", func(t *testing.T) {
		app := &api{reviews: &fakeReviewUseCase{
			listForServiceFunc: func(ctx context.Context, serviceID int) ([]Review, error) {
				return []Review{{ID: 1, ExchangeID: 3, Rating: 4}}, nil
			},
		}}

		req := httptest.NewRequest(http.MethodGet, "/api/services/1/reviews", nil)
		req.SetPathValue("id", "1")
		rec := httptest.NewRecorder()

		app.handleListServiceReviews(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200", rec.Code)
		}
	})
}
