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

	t.Run("mise a jour de son profil -> 200", func(t *testing.T) {
		app := &api{users: &fakeUserUseCase{
			updateProfileFunc: func(ctx context.Context, actorID, targetID int, pseudo, bio, ville string) (User, error) {
				return User{ID: targetID, Pseudo: pseudo, Bio: bio, Ville: ville, CreditBalance: 10}, nil
			},
		}}

		body := `{"pseudo":"Thierry","bio":"bio","ville":"Lyon"}`
		req := httptest.NewRequest(http.MethodPut, "/api/users/5", strings.NewReader(body))
		req.SetPathValue("id", "5")
		req = withCurrentUser(req, User{ID: 5})
		rec := httptest.NewRecorder()

		app.handleUpdateUser(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("code = %d, attendu %d", rec.Code, http.StatusOK)
		}
		var got User
		if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
			t.Fatalf("reponse JSON invalide : %v", err)
		}
		if got.Pseudo != "Thierry" {
			t.Errorf("Pseudo = %q, attendu Thierry", got.Pseudo)
		}
		if got.Ville != "Lyon" {
			t.Errorf("Ville = %q, attendu Lyon", got.Ville)
		}
	})

	t.Run("profil d'un autre utilisateur -> 403", func(t *testing.T) {
		app := &api{users: &fakeUserUseCase{
			updateProfileFunc: func(ctx context.Context, actorID, targetID int, pseudo, bio, ville string) (User, error) {
				return User{}, ErrForbidden
			},
		}}

		req := httptest.NewRequest(http.MethodPut, "/api/users/9", strings.NewReader(`{"pseudo":"Thierry"}`))
		req.SetPathValue("id", "9")
		req = withCurrentUser(req, User{ID: 5})
		rec := httptest.NewRecorder()

		app.handleUpdateUser(rec, req)

		if rec.Code != http.StatusForbidden {
			t.Fatalf("code = %d, attendu %d", rec.Code, http.StatusForbidden)
		}
	})

	t.Run("pseudo vide -> 400", func(t *testing.T) {
		app := &api{users: &fakeUserUseCase{
			updateProfileFunc: func(ctx context.Context, actorID, targetID int, pseudo, bio, ville string) (User, error) {
				return User{}, ErrPseudoRequired
			},
		}}

		req := httptest.NewRequest(http.MethodPut, "/api/users/5", strings.NewReader(`{"pseudo":""}`))
		req.SetPathValue("id", "5")
		req = withCurrentUser(req, User{ID: 5})
		rec := httptest.NewRecorder()

		app.handleUpdateUser(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("code = %d, attendu %d", rec.Code, http.StatusBadRequest)
		}
	})

	t.Run("sans utilisateur authentifie -> 401", func(t *testing.T) {
		app := &api{users: &fakeUserUseCase{}}

		req := httptest.NewRequest(http.MethodPut, "/api/users/5", strings.NewReader(`{"pseudo":"Thierry"}`))
		req.SetPathValue("id", "5")
		rec := httptest.NewRecorder()

		app.handleUpdateUser(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("code = %d, attendu %d", rec.Code, http.StatusUnauthorized)
		}
	})

	t.Run("identifiant invalide -> 400", func(t *testing.T) {
		app := &api{users: &fakeUserUseCase{}}

		req := httptest.NewRequest(http.MethodPut, "/api/users/abc", strings.NewReader(`{"pseudo":"Thierry"}`))
		req.SetPathValue("id", "abc")
		req = withCurrentUser(req, User{ID: 5})
		rec := httptest.NewRecorder()

		app.handleUpdateUser(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("code = %d, attendu %d", rec.Code, http.StatusBadRequest)
		}
	})
}
