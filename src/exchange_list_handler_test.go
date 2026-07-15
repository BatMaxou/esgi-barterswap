package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleListExchanges(t *testing.T) {
	t.Run("status query param forwarded -> 200", func(t *testing.T) {
		var capturedStatus string
		app := &api{exchanges: &fakeExchangeUseCase{
			listFunc: func(ctx context.Context, actorID int, status string) ([]Exchange, error) {
				capturedStatus = status
				return []Exchange{{ID: 1, Status: ExchangeStatusPending}}, nil
			},
		}}

		req := httptest.NewRequest(http.MethodGet, "/api/exchanges?status=pending", nil)
		req = withUser(req, User{ID: 2})
		rec := httptest.NewRecorder()

		app.handleListExchanges(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200", rec.Code)
		}
		if capturedStatus != ExchangeStatusPending {
			t.Errorf("status filter = %q, want pending", capturedStatus)
		}
	})

	t.Run("without authenticated user -> 401", func(t *testing.T) {
		app := &api{exchanges: &fakeExchangeUseCase{}}

		req := httptest.NewRequest(http.MethodGet, "/api/exchanges", nil)
		rec := httptest.NewRecorder()

		app.handleListExchanges(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("status = %d, want 401", rec.Code)
		}
	})
}
