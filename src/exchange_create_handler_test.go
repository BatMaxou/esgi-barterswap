package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type fakeExchangeUseCase struct {
	createFunc   func(ctx context.Context, requesterID, serviceID int) (Exchange, error)
	listFunc     func(ctx context.Context, actorID int, status string) ([]Exchange, error)
	getFunc      func(ctx context.Context, actorID, exchangeID int) (Exchange, error)
	acceptFunc   func(ctx context.Context, actorID, exchangeID int) (Exchange, error)
	rejectFunc   func(ctx context.Context, actorID, exchangeID int) (Exchange, error)
	completeFunc func(ctx context.Context, actorID, exchangeID int) (Exchange, error)
	cancelFunc   func(ctx context.Context, actorID, exchangeID int) (Exchange, error)
}

func (fake *fakeExchangeUseCase) Create(ctx context.Context, requesterID, serviceID int) (Exchange, error) {
	return fake.createFunc(ctx, requesterID, serviceID)
}

func (fake *fakeExchangeUseCase) List(ctx context.Context, actorID int, status string) ([]Exchange, error) {
	return fake.listFunc(ctx, actorID, status)
}

func (fake *fakeExchangeUseCase) Get(ctx context.Context, actorID, exchangeID int) (Exchange, error) {
	return fake.getFunc(ctx, actorID, exchangeID)
}

func (fake *fakeExchangeUseCase) Accept(ctx context.Context, actorID, exchangeID int) (Exchange, error) {
	return fake.acceptFunc(ctx, actorID, exchangeID)
}

func (fake *fakeExchangeUseCase) Reject(ctx context.Context, actorID, exchangeID int) (Exchange, error) {
	return fake.rejectFunc(ctx, actorID, exchangeID)
}

func (fake *fakeExchangeUseCase) Complete(ctx context.Context, actorID, exchangeID int) (Exchange, error) {
	return fake.completeFunc(ctx, actorID, exchangeID)
}

func (fake *fakeExchangeUseCase) Cancel(ctx context.Context, actorID, exchangeID int) (Exchange, error) {
	return fake.cancelFunc(ctx, actorID, exchangeID)
}

func TestHandleCreateExchange(t *testing.T) {
	t.Run("valid request -> 201", func(t *testing.T) {
		app := &api{exchanges: &fakeExchangeUseCase{
			createFunc: func(ctx context.Context, requesterID, serviceID int) (Exchange, error) {
				return Exchange{ID: 1, ServiceID: serviceID, RequesterID: requesterID, Status: ExchangeStatusPending}, nil
			},
		}}

		body := `{"service_id":3}`
		req := httptest.NewRequest(http.MethodPost, "/api/exchanges", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req = withUser(req, User{ID: 2})
		rec := httptest.NewRecorder()

		app.handleCreateExchange(rec, req)

		if rec.Code != http.StatusCreated {
			t.Fatalf("status = %d, want 201", rec.Code)
		}
	})

	t.Run("self request -> 400", func(t *testing.T) {
		app := &api{exchanges: &fakeExchangeUseCase{
			createFunc: func(ctx context.Context, requesterID, serviceID int) (Exchange, error) {
				return Exchange{}, ErrExchangeSelfRequest
			},
		}}

		req := httptest.NewRequest(http.MethodPost, "/api/exchanges", strings.NewReader(`{"service_id":3}`))
		req.Header.Set("Content-Type", "application/json")
		req = withUser(req, User{ID: 2})
		rec := httptest.NewRecorder()

		app.handleCreateExchange(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want 400", rec.Code)
		}
	})

	t.Run("service unavailable -> 409", func(t *testing.T) {
		app := &api{exchanges: &fakeExchangeUseCase{
			createFunc: func(ctx context.Context, requesterID, serviceID int) (Exchange, error) {
				return Exchange{}, ErrExchangeServiceUnavailable
			},
		}}

		req := httptest.NewRequest(http.MethodPost, "/api/exchanges", strings.NewReader(`{"service_id":3}`))
		req.Header.Set("Content-Type", "application/json")
		req = withUser(req, User{ID: 2})
		rec := httptest.NewRecorder()

		app.handleCreateExchange(rec, req)

		if rec.Code != http.StatusConflict {
			t.Fatalf("status = %d, want 409", rec.Code)
		}
	})
}
