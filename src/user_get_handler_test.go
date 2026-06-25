package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleGetUser(t *testing.T) {
	t.Run("utilisateur existant -> 200", func(t *testing.T) {
		app := &api{users: &fakeUserUseCase{
			getProfileFunc: func(ctx context.Context, id int) (User, error) {
				return User{
					ID:            id,
					Pseudo:        "Thierry",
					Skills:        []Skill{},
					CreditBalance: 35,
					CreatedAt:     "2026-01-01T00:00:00Z",
				}, nil
			},
		}}

		req := httptest.NewRequest(http.MethodGet, "/api/users/5", nil)
		req.SetPathValue("id", "5")
		rec := httptest.NewRecorder()

		app.handleGetUser(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("code = %d, attendu %d", rec.Code, http.StatusOK)
		}

		var got User
		if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
			t.Fatalf("reponse JSON invalide : %v", err)
		}
		if got.ID != 5 {
			t.Errorf("ID = %d, attendu 5", got.ID)
		}
		if got.CreditBalance != 35 {
			t.Errorf("CreditBalance = %d, attendu 35", got.CreditBalance)
		}
	})

	t.Run("utilisateur introuvable -> 404", func(t *testing.T) {
		app := &api{users: &fakeUserUseCase{
			getProfileFunc: func(ctx context.Context, id int) (User, error) {
				return User{}, ErrUserNotFound
			},
		}}

		req := httptest.NewRequest(http.MethodGet, "/api/users/999", nil)
		req.SetPathValue("id", "999")
		rec := httptest.NewRecorder()

		app.handleGetUser(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Fatalf("code = %d, attendu %d", rec.Code, http.StatusNotFound)
		}
	})

	t.Run("identifiant invalide -> 400", func(t *testing.T) {
		app := &api{users: &fakeUserUseCase{}}

		req := httptest.NewRequest(http.MethodGet, "/api/users/abc", nil)
		req.SetPathValue("id", "abc")
		rec := httptest.NewRecorder()

		app.handleGetUser(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("code = %d, attendu %d", rec.Code, http.StatusBadRequest)
		}
	})
}
