package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequireAuth(t *testing.T) {
	newApp := func(authenticate func(ctx context.Context, id int) (User, error)) *api {
		return &api{users: &fakeUserUseCase{authenticateFunc: authenticate}}
	}

	t.Run("missing X-User-ID header -> 401", func(t *testing.T) {
		called := false
		app := newApp(func(ctx context.Context, id int) (User, error) {
			t.Fatal("Authenticate must not be called without a header")

			return User{}, nil
		})
		next := func(w http.ResponseWriter, r *http.Request) { called = true }

		req := httptest.NewRequest(http.MethodPut, "/api/users/5", nil)
		rec := httptest.NewRecorder()
		app.requireAuth(next)(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
		}
		if called {
			t.Error("the next handler must not be called")
		}
	})

	t.Run("non-integer X-User-ID -> 401", func(t *testing.T) {
		called := false
		app := newApp(func(ctx context.Context, id int) (User, error) {
			t.Fatal("Authenticate must not be called for a non-integer id")

			return User{}, nil
		})
		next := func(w http.ResponseWriter, r *http.Request) { called = true }

		req := httptest.NewRequest(http.MethodPut, "/api/users/5", nil)
		req.Header.Set("X-User-ID", "abc")
		rec := httptest.NewRecorder()
		app.requireAuth(next)(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
		}
		if called {
			t.Error("the next handler must not be called")
		}
	})

	t.Run("unknown user -> 401", func(t *testing.T) {
		called := false
		app := newApp(func(ctx context.Context, id int) (User, error) {
			return User{}, ErrUserNotFound
		})
		next := func(w http.ResponseWriter, r *http.Request) { called = true }

		req := httptest.NewRequest(http.MethodPut, "/api/users/999", nil)
		req.Header.Set("X-User-ID", "999")
		rec := httptest.NewRecorder()
		app.requireAuth(next)(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
		}
		if called {
			t.Error("the next handler must not be called")
		}
	})

	t.Run("valid user -> handler called with the current user", func(t *testing.T) {
		called := false
		var seenUser User
		app := newApp(func(ctx context.Context, id int) (User, error) {
			return User{ID: id, Pseudo: "Thierry"}, nil
		})
		next := func(w http.ResponseWriter, r *http.Request) {
			user, ok := currentUser(r.Context())
			if !ok {
				t.Fatal("current user missing from context")
			}
			called = true
			seenUser = user
			w.WriteHeader(http.StatusOK)
		}

		req := httptest.NewRequest(http.MethodPut, "/api/users/5", nil)
		req.Header.Set("X-User-ID", "5")
		rec := httptest.NewRecorder()
		app.requireAuth(next)(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
		}
		if !called {
			t.Fatal("the next handler must be called")
		}
		if seenUser.ID != 5 {
			t.Errorf("current user ID = %d, want 5", seenUser.ID)
		}
	})
}
