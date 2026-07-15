package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleCancelExchange(t *testing.T) {
	t.Run("participant cancels -> 200", func(t *testing.T) {
		app := &api{exchanges: &fakeExchangeUseCase{
			cancelFunc: func(ctx context.Context, actorID, exchangeID int) (Exchange, error) {
				return Exchange{ID: exchangeID, Status: ExchangeStatusCancelled}, nil
			},
		}}

		req := httptest.NewRequest(http.MethodPut, "/api/exchanges/5/cancel", nil)
		req.SetPathValue("id", "5")
		req = withUser(req, User{ID: 2})
		rec := httptest.NewRecorder()

		app.handleCancelExchange(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200", rec.Code)
		}
	})
}
