package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandleUpdateService(t *testing.T) {
	t.Run("updating your own ad -> 200", func(t *testing.T) {
		app := &api{services: &fakeServiceUseCase{
			updateFunc: func(ctx context.Context, actorID, serviceID int, title, description, category, city string, durationMinutes, credits int, active *bool) (Service, error) {
				return Service{ID: serviceID, ProviderID: actorID, Title: title}, nil
			},
		}}

		body := `{"title":"Nouveau titre","category":"Informatique","duration_minutes":90,"credits":3}`
		req := httptest.NewRequest(http.MethodPut, "/api/services/7", strings.NewReader(body))
		req.SetPathValue("id", "7")
		req = withUser(req, User{ID: 5})
		rec := httptest.NewRecorder()

		app.handleUpdateService(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
		}
	})

	t.Run("another user's ad -> 403", func(t *testing.T) {
		app := &api{services: &fakeServiceUseCase{
			updateFunc: func(ctx context.Context, actorID, serviceID int, title, description, category, city string, durationMinutes, credits int, active *bool) (Service, error) {
				return Service{}, ErrForbidden
			},
		}}

		req := httptest.NewRequest(http.MethodPut, "/api/services/7", strings.NewReader(`{"title":"X","category":"Informatique","duration_minutes":90,"credits":3}`))
		req.SetPathValue("id", "7")
		req = withUser(req, User{ID: 9})
		rec := httptest.NewRecorder()

		app.handleUpdateService(rec, req)

		if rec.Code != http.StatusForbidden {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusForbidden)
		}
	})

	t.Run("ad not found -> 404", func(t *testing.T) {
		app := &api{services: &fakeServiceUseCase{
			updateFunc: func(ctx context.Context, actorID, serviceID int, title, description, category, city string, durationMinutes, credits int, active *bool) (Service, error) {
				return Service{}, ErrServiceNotFound
			},
		}}

		req := httptest.NewRequest(http.MethodPut, "/api/services/999", strings.NewReader(`{"title":"X","category":"Informatique","duration_minutes":90,"credits":3}`))
		req.SetPathValue("id", "999")
		req = withUser(req, User{ID: 5})
		rec := httptest.NewRecorder()

		app.handleUpdateService(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
		}
	})

	t.Run("validation failed -> 400", func(t *testing.T) {
		app := &api{services: &fakeServiceUseCase{
			updateFunc: func(ctx context.Context, actorID, serviceID int, title, description, category, city string, durationMinutes, credits int, active *bool) (Service, error) {
				return Service{}, ErrServiceDurationInvalid
			},
		}}

		req := httptest.NewRequest(http.MethodPut, "/api/services/7", strings.NewReader(`{"title":"X","category":"Informatique","duration_minutes":0,"credits":3}`))
		req.SetPathValue("id", "7")
		req = withUser(req, User{ID: 5})
		rec := httptest.NewRecorder()

		app.handleUpdateService(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
	})

	t.Run("without authenticated user -> 401", func(t *testing.T) {
		app := &api{services: &fakeServiceUseCase{}}

		req := httptest.NewRequest(http.MethodPut, "/api/services/7", strings.NewReader(`{}`))
		req.SetPathValue("id", "7")
		rec := httptest.NewRecorder()

		app.handleUpdateService(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
		}
	})
}
