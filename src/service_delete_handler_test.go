package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleDeleteService(t *testing.T) {
	t.Run("deleting your own ad -> 204", func(t *testing.T) {
		app := &api{services: &fakeServiceUseCase{
			deleteFunc: func(ctx context.Context, actorID, serviceID int) error {
				return nil
			},
		}}

		req := httptest.NewRequest(http.MethodDelete, "/api/services/7", nil)
		req.SetPathValue("id", "7")
		req = withUser(req, User{ID: 5})
		rec := httptest.NewRecorder()

		app.handleDeleteService(rec, req)

		if rec.Code != http.StatusNoContent {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusNoContent)
		}
	})

	t.Run("another user's ad -> 403", func(t *testing.T) {
		app := &api{services: &fakeServiceUseCase{
			deleteFunc: func(ctx context.Context, actorID, serviceID int) error {
				return ErrForbidden
			},
		}}

		req := httptest.NewRequest(http.MethodDelete, "/api/services/7", nil)
		req.SetPathValue("id", "7")
		req = withUser(req, User{ID: 9})
		rec := httptest.NewRecorder()

		app.handleDeleteService(rec, req)

		if rec.Code != http.StatusForbidden {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusForbidden)
		}
	})

	t.Run("ad not found -> 404", func(t *testing.T) {
		app := &api{services: &fakeServiceUseCase{
			deleteFunc: func(ctx context.Context, actorID, serviceID int) error {
				return ErrServiceNotFound
			},
		}}

		req := httptest.NewRequest(http.MethodDelete, "/api/services/999", nil)
		req.SetPathValue("id", "999")
		req = withUser(req, User{ID: 5})
		rec := httptest.NewRecorder()

		app.handleDeleteService(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
		}
	})

	t.Run("ad referenced by an exchange -> 409", func(t *testing.T) {
		app := &api{services: &fakeServiceUseCase{
			deleteFunc: func(ctx context.Context, actorID, serviceID int) error {
				return ErrServiceHasExchanges
			},
		}}

		req := httptest.NewRequest(http.MethodDelete, "/api/services/7", nil)
		req.SetPathValue("id", "7")
		req = withUser(req, User{ID: 5})
		rec := httptest.NewRecorder()

		app.handleDeleteService(rec, req)

		if rec.Code != http.StatusConflict {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusConflict)
		}
	})

	t.Run("without authenticated user -> 401", func(t *testing.T) {
		app := &api{services: &fakeServiceUseCase{}}

		req := httptest.NewRequest(http.MethodDelete, "/api/services/7", nil)
		req.SetPathValue("id", "7")
		rec := httptest.NewRecorder()

		app.handleDeleteService(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
		}
	})
}
