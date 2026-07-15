package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type fakeServiceUseCase struct {
	createFunc func(ctx context.Context, providerID int, title, description, category, city string, durationMinutes, credits int) (Service, error)
	getFunc    func(ctx context.Context, id int) (Service, error)
	listFunc   func(ctx context.Context, filter ServiceFilter) ([]Service, error)
	updateFunc func(ctx context.Context, actorID, serviceID int, title, description, category, city string, durationMinutes, credits int, active *bool) (Service, error)
	deleteFunc func(ctx context.Context, actorID, serviceID int) error
}

func (fake *fakeServiceUseCase) Create(ctx context.Context, providerID int, title, description, category, city string, durationMinutes, credits int) (Service, error) {
	return fake.createFunc(ctx, providerID, title, description, category, city, durationMinutes, credits)
}

func (fake *fakeServiceUseCase) Get(ctx context.Context, id int) (Service, error) {
	return fake.getFunc(ctx, id)
}

func (fake *fakeServiceUseCase) List(ctx context.Context, filter ServiceFilter) ([]Service, error) {
	return fake.listFunc(ctx, filter)
}

func (fake *fakeServiceUseCase) Update(ctx context.Context, actorID, serviceID int, title, description, category, city string, durationMinutes, credits int, active *bool) (Service, error) {
	return fake.updateFunc(ctx, actorID, serviceID, title, description, category, city, durationMinutes, credits, active)
}

func (fake *fakeServiceUseCase) Delete(ctx context.Context, actorID, serviceID int) error {
	return fake.deleteFunc(ctx, actorID, serviceID)
}

// withUser attaches an authenticated user to the request context, like the
// requireAuth middleware does.
func withUser(req *http.Request, user User) *http.Request {
	ctx := context.WithValue(req.Context(), currentUserKey, user)
	return req.WithContext(ctx)
}

func TestHandleCreateService(t *testing.T) {
	t.Run("valid creation -> 201", func(t *testing.T) {
		app := &api{services: &fakeServiceUseCase{
			createFunc: func(ctx context.Context, providerID int, title, description, category, city string, durationMinutes, credits int) (Service, error) {
				return Service{ID: 42, ProviderID: providerID, Title: title, Category: category, Active: true}, nil
			},
		}}

		body := `{"title":"Cours de Go","category":"Informatique","duration_minutes":60,"credits":2}`
		req := httptest.NewRequest(http.MethodPost, "/api/services", strings.NewReader(body))
		req = withUser(req, User{ID: 5})
		rec := httptest.NewRecorder()

		app.handleCreateService(rec, req)

		if rec.Code != http.StatusCreated {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusCreated)
		}
		var got Service
		if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
			t.Fatalf("invalid JSON response: %v", err)
		}
		if got.ProviderID != 5 {
			t.Errorf("ProviderID = %d, want 5 (authenticated user)", got.ProviderID)
		}
	})

	t.Run("invalid category -> 400", func(t *testing.T) {
		app := &api{services: &fakeServiceUseCase{
			createFunc: func(ctx context.Context, providerID int, title, description, category, city string, durationMinutes, credits int) (Service, error) {
				return Service{}, ErrServiceCategoryInvalid
			},
		}}

		req := httptest.NewRequest(http.MethodPost, "/api/services", strings.NewReader(`{"title":"X","category":"Truc","duration_minutes":60,"credits":2}`))
		req = withUser(req, User{ID: 5})
		rec := httptest.NewRecorder()

		app.handleCreateService(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
	})

	t.Run("without authenticated user -> 401", func(t *testing.T) {
		app := &api{services: &fakeServiceUseCase{}}

		req := httptest.NewRequest(http.MethodPost, "/api/services", strings.NewReader(`{}`))
		rec := httptest.NewRecorder()

		app.handleCreateService(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
		}
	})

	t.Run("invalid JSON -> 400", func(t *testing.T) {
		app := &api{services: &fakeServiceUseCase{}}

		req := httptest.NewRequest(http.MethodPost, "/api/services", strings.NewReader(`{not json`))
		req = withUser(req, User{ID: 5})
		rec := httptest.NewRecorder()

		app.handleCreateService(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
	})
}
