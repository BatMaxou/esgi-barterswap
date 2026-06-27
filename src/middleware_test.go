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

	t.Run("header X-User-ID absent -> 401", func(t *testing.T) {
		called := false
		app := newApp(func(ctx context.Context, id int) (User, error) {
			t.Fatal("Authenticate ne doit pas etre appele sans header")

			return User{}, nil
		})
		next := func(w http.ResponseWriter, r *http.Request) { called = true }

		req := httptest.NewRequest(http.MethodPut, "/api/users/5", nil)
		rec := httptest.NewRecorder()
		app.requireAuth(next)(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("code = %d, attendu %d", rec.Code, http.StatusUnauthorized)
		}
		if called {
			t.Error("le handler suivant ne doit pas etre appele")
		}
	})

	t.Run("X-User-ID non entier -> 401", func(t *testing.T) {
		called := false
		app := newApp(func(ctx context.Context, id int) (User, error) {
			t.Fatal("Authenticate ne doit pas etre appele pour un id non entier")

			return User{}, nil
		})
		next := func(w http.ResponseWriter, r *http.Request) { called = true }

		req := httptest.NewRequest(http.MethodPut, "/api/users/5", nil)
		req.Header.Set("X-User-ID", "abc")
		rec := httptest.NewRecorder()
		app.requireAuth(next)(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("code = %d, attendu %d", rec.Code, http.StatusUnauthorized)
		}
		if called {
			t.Error("le handler suivant ne doit pas etre appele")
		}
	})

	t.Run("utilisateur inconnu -> 401", func(t *testing.T) {
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
			t.Fatalf("code = %d, attendu %d", rec.Code, http.StatusUnauthorized)
		}
		if called {
			t.Error("le handler suivant ne doit pas etre appele")
		}
	})

	t.Run("utilisateur valide -> handler appele avec l'utilisateur courant", func(t *testing.T) {
		called := false
		var seenUser User
		app := newApp(func(ctx context.Context, id int) (User, error) {
			return User{ID: id, Pseudo: "Thierry"}, nil
		})
		next := func(w http.ResponseWriter, r *http.Request) {
			user, ok := currentUser(r.Context())
			if !ok {
				t.Fatal("utilisateur courant absent du context")
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
			t.Fatalf("code = %d, attendu %d", rec.Code, http.StatusOK)
		}
		if !called {
			t.Fatal("le handler suivant doit etre appele")
		}
		if seenUser.ID != 5 {
			t.Errorf("utilisateur courant ID = %d, attendu 5", seenUser.ID)
		}
	})
}
