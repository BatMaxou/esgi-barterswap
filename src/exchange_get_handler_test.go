package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleGetExchange(t *testing.T) {
	t.Run("participant -> 200", func(t *testing.T) {
		app := &api{exchanges: &fakeExchangeUseCase{
			getFunc: func(ctx context.Context, actorID, exchangeID int) (Exchange, error) {
				return Exchange{ID: exchangeID, RequesterID: actorID}, nil
			},
		}}

		req := httptest.NewRequest(http.MethodGet, "/api/exchanges/4", nil)
		req.SetPathValue("id", "4")
		req = withUser(req, User{ID: 2})
		rec := httptest.NewRecorder()

		app.handleGetExchange(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200", rec.Code)
		}
	})

	t.Run("not found -> 404", func(t *testing.T) {
		app := &api{exchanges: &fakeExchangeUseCase{
			getFunc: func(ctx context.Context, actorID, exchangeID int) (Exchange, error) {
				return Exchange{}, ErrExchangeNotFound
			},
		}}

		req := httptest.NewRequest(http.MethodGet, "/api/exchanges/999", nil)
		req.SetPathValue("id", "999")
		req = withUser(req, User{ID: 2})
		rec := httptest.NewRecorder()

		app.handleGetExchange(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want 404", rec.Code)
		}
	})

	t.Run("invalid id -> 400", func(t *testing.T) {
		app := &api{exchanges: &fakeExchangeUseCase{}}

		req := httptest.NewRequest(http.MethodGet, "/api/exchanges/abc", nil)
		req.SetPathValue("id", "abc")
		req = withUser(req, User{ID: 2})
		rec := httptest.NewRecorder()

		app.handleGetExchange(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want 400", rec.Code)
		}
	})
}
