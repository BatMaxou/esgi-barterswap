package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandleUpdateUser(t *testing.T) {
	withCurrentUser := func(req *http.Request, user User) *http.Request {
		ctx := context.WithValue(req.Context(), currentUserKey, user)
		return req.WithContext(ctx)
	}

	t.Run("updating your own profile -> 200", func(t *testing.T) {
		app := &api{users: &fakeUserUseCase{
			updateProfileFunc: func(ctx context.Context, actorID, targetID int, pseudo, bio, city string) (User, error) {
				return User{ID: targetID, Pseudo: pseudo, Bio: bio, City: city, CreditBalance: 10}, nil
			},
		}}

		body := `{"pseudo":"Thierry","bio":"bio","city":"Lyon"}`
		req := httptest.NewRequest(http.MethodPut, "/api/users/5", strings.NewReader(body))
		req.SetPathValue("id", "5")
		req = withCurrentUser(req, User{ID: 5})
		rec := httptest.NewRecorder()

		app.handleUpdateUser(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
		}
		var got User
		if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
			t.Fatalf("invalid JSON response: %v", err)
		}
		if got.Pseudo != "Thierry" {
			t.Errorf("Pseudo = %q, want Thierry", got.Pseudo)
		}
		if got.City != "Lyon" {
			t.Errorf("City = %q, want Lyon", got.City)
		}
	})

	t.Run("another user's profile -> 403", func(t *testing.T) {
		app := &api{users: &fakeUserUseCase{
			updateProfileFunc: func(ctx context.Context, actorID, targetID int, pseudo, bio, city string) (User, error) {
				return User{}, ErrForbidden
			},
		}}

		req := httptest.NewRequest(http.MethodPut, "/api/users/9", strings.NewReader(`{"pseudo":"Thierry"}`))
		req.SetPathValue("id", "9")
		req = withCurrentUser(req, User{ID: 5})
		rec := httptest.NewRecorder()

		app.handleUpdateUser(rec, req)

		if rec.Code != http.StatusForbidden {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusForbidden)
		}
	})

	t.Run("empty pseudo -> 400", func(t *testing.T) {
		app := &api{users: &fakeUserUseCase{
			updateProfileFunc: func(ctx context.Context, actorID, targetID int, pseudo, bio, city string) (User, error) {
				return User{}, ErrPseudoRequired
			},
		}}

		req := httptest.NewRequest(http.MethodPut, "/api/users/5", strings.NewReader(`{"pseudo":""}`))
		req.SetPathValue("id", "5")
		req = withCurrentUser(req, User{ID: 5})
		rec := httptest.NewRecorder()

		app.handleUpdateUser(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
	})

	t.Run("without authenticated user -> 401", func(t *testing.T) {
		app := &api{users: &fakeUserUseCase{}}

		req := httptest.NewRequest(http.MethodPut, "/api/users/5", strings.NewReader(`{"pseudo":"Thierry"}`))
		req.SetPathValue("id", "5")
		rec := httptest.NewRecorder()

		app.handleUpdateUser(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
		}
	})

	t.Run("invalid identifier -> 400", func(t *testing.T) {
		app := &api{users: &fakeUserUseCase{}}

		req := httptest.NewRequest(http.MethodPut, "/api/users/abc", strings.NewReader(`{"pseudo":"Thierry"}`))
		req.SetPathValue("id", "abc")
		req = withCurrentUser(req, User{ID: 5})
		rec := httptest.NewRecorder()

		app.handleUpdateUser(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
	})
}
