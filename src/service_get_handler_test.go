package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleGetService(t *testing.T) {
	t.Run("existing service -> 200", func(t *testing.T) {
		app := &api{services: &fakeServiceUseCase{
			getFunc: func(ctx context.Context, id int) (Service, error) {
				return Service{ID: id, Title: "Cours de Go"}, nil
			},
		}}

		req := httptest.NewRequest(http.MethodGet, "/api/services/7", nil)
		req.SetPathValue("id", "7")
		rec := httptest.NewRecorder()

		app.handleGetService(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
		}
	})

	t.Run("service not found -> 404", func(t *testing.T) {
		app := &api{services: &fakeServiceUseCase{
			getFunc: func(ctx context.Context, id int) (Service, error) {
				return Service{}, ErrServiceNotFound
			},
		}}

		req := httptest.NewRequest(http.MethodGet, "/api/services/999", nil)
		req.SetPathValue("id", "999")
		rec := httptest.NewRecorder()

		app.handleGetService(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
		}
	})

	t.Run("non-numeric identifier -> 400", func(t *testing.T) {
		app := &api{services: &fakeServiceUseCase{}}

		req := httptest.NewRequest(http.MethodGet, "/api/services/abc", nil)
		req.SetPathValue("id", "abc")
		rec := httptest.NewRecorder()

		app.handleGetService(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
	})
}

func TestHandleListServices(t *testing.T) {
	t.Run("filters read from the query params -> 200", func(t *testing.T) {
		var captured ServiceFilter
		app := &api{services: &fakeServiceUseCase{
			listFunc: func(ctx context.Context, filter ServiceFilter) ([]Service, error) {
				captured = filter
				return []Service{{ID: 1, Title: "Cours de Go"}}, nil
			},
		}}

		req := httptest.NewRequest(http.MethodGet, "/api/services?category=Informatique&city=Paris&search=Go", nil)
		rec := httptest.NewRecorder()

		app.handleListServices(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
		}
		if captured.Category != "Informatique" || captured.City != "Paris" || captured.Search != "Go" {
			t.Errorf("filter = %+v, want Informatique/Paris/Go", captured)
		}
	})
}
