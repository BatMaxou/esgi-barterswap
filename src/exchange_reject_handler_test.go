package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleRejectExchange(t *testing.T) {
	t.Run("owner rejects -> 200", func(t *testing.T) {
		app := &api{exchanges: &fakeExchangeUseCase{
			rejectFunc: func(ctx context.Context, actorID, exchangeID int) (Exchange, error) {
				return Exchange{ID: exchangeID, Status: ExchangeStatusRejected}, nil
			},
		}}

		req := httptest.NewRequest(http.MethodPut, "/api/exchanges/5/reject", nil)
		req.SetPathValue("id", "5")
		req = withUser(req, User{ID: 1})
		rec := httptest.NewRecorder()

		app.handleRejectExchange(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200", rec.Code)
		}
	})

	t.Run("invalid transition -> 400", func(t *testing.T) {
		app := &api{exchanges: &fakeExchangeUseCase{
			rejectFunc: func(ctx context.Context, actorID, exchangeID int) (Exchange, error) {
				return Exchange{}, ErrExchangeInvalidTransition
			},
		}}

		req := httptest.NewRequest(http.MethodPut, "/api/exchanges/5/reject", nil)
		req.SetPathValue("id", "5")
		req = withUser(req, User{ID: 1})
		rec := httptest.NewRecorder()

		app.handleRejectExchange(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want 400", rec.Code)
		}
	})
}
