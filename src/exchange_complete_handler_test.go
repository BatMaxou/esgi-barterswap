package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleCompleteExchange(t *testing.T) {
	t.Run("owner completes -> 200", func(t *testing.T) {
		app := &api{exchanges: &fakeExchangeUseCase{
			completeFunc: func(ctx context.Context, actorID, exchangeID int) (Exchange, error) {
				return Exchange{ID: exchangeID, Status: ExchangeStatusCompleted}, nil
			},
		}}

		req := httptest.NewRequest(http.MethodPut, "/api/exchanges/5/complete", nil)
		req.SetPathValue("id", "5")
		req = withUser(req, User{ID: 1})
		rec := httptest.NewRecorder()

		app.handleCompleteExchange(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200", rec.Code)
		}
	})
}
