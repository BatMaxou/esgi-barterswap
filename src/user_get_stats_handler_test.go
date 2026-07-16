package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

type fakeUserStatsUseCase struct {
	getFunc func(ctx context.Context, userID int) (UserStats, error)
}

func (fake *fakeUserStatsUseCase) Get(ctx context.Context, userID int) (UserStats, error) {
	return fake.getFunc(ctx, userID)
}

func TestHandleGetUserStats(t *testing.T) {
	t.Run("existing user -> 200", func(t *testing.T) {
		app := &api{stats: &fakeUserStatsUseCase{
			getFunc: func(ctx context.Context, userID int) (UserStats, error) {
				return UserStats{UserID: userID, CreditBalance: 10, TotalEarned: 10}, nil
			},
		}}

		req := httptest.NewRequest(http.MethodGet, "/api/users/1/stats", nil)
		req.SetPathValue("id", "1")
		rec := httptest.NewRecorder()

		app.handleGetUserStats(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200", rec.Code)
		}
	})

	t.Run("user not found -> 404", func(t *testing.T) {
		app := &api{stats: &fakeUserStatsUseCase{
			getFunc: func(ctx context.Context, userID int) (UserStats, error) {
				return UserStats{}, ErrUserNotFound
			},
		}}

		req := httptest.NewRequest(http.MethodGet, "/api/users/999/stats", nil)
		req.SetPathValue("id", "999")
		rec := httptest.NewRecorder()

		app.handleGetUserStats(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want 404", rec.Code)
		}
	})

	t.Run("invalid id -> 400", func(t *testing.T) {
		app := &api{stats: &fakeUserStatsUseCase{}}

		req := httptest.NewRequest(http.MethodGet, "/api/users/abc/stats", nil)
		req.SetPathValue("id", "abc")
		rec := httptest.NewRecorder()

		app.handleGetUserStats(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want 400", rec.Code)
		}
	})
}
