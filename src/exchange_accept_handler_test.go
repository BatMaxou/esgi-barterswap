package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleAcceptExchange(t *testing.T) {
	t.Run("owner accepts -> 200", func(t *testing.T) {
		app := &api{exchanges: &fakeExchangeUseCase{
			acceptFunc: func(ctx context.Context, actorID, exchangeID int) (Exchange, error) {
				return Exchange{ID: exchangeID, OwnerID: actorID, Status: ExchangeStatusAccepted}, nil
			},
		}}

		req := httptest.NewRequest(http.MethodPut, "/api/exchanges/5/accept", nil)
		req.SetPathValue("id", "5")
		req = withUser(req, User{ID: 1})
		rec := httptest.NewRecorder()

		app.handleAcceptExchange(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200", rec.Code)
		}
	})

	t.Run("forbidden -> 403", func(t *testing.T) {
		app := &api{exchanges: &fakeExchangeUseCase{
			acceptFunc: func(ctx context.Context, actorID, exchangeID int) (Exchange, error) {
				return Exchange{}, ErrForbidden
			},
		}}

		req := httptest.NewRequest(http.MethodPut, "/api/exchanges/5/accept", nil)
		req.SetPathValue("id", "5")
		req = withUser(req, User{ID: 2})
		rec := httptest.NewRecorder()

		app.handleAcceptExchange(rec, req)

		if rec.Code != http.StatusForbidden {
			t.Fatalf("status = %d, want 403", rec.Code)
		}
	})
}
