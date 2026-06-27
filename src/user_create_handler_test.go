package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandleCreateUser(t *testing.T) {
	t.Run("creation reussie -> 201", func(t *testing.T) {
		app := &api{users: &fakeUserUseCase{
			registerFunc: func(ctx context.Context, pseudo, bio, city string) (User, error) {
				return User{
					ID:            42,
					Pseudo:        pseudo,
					Bio:           bio,
					City:          city,
					Skills:        []Skill{},
					CreditBalance: welcomeCredits,
					CreatedAt:     "2026-01-01T00:00:00Z",
				}, nil
			},
		}}

		body := `{"pseudo":"Thierry","bio":"ma bio","city":"Paris"}`
		req := httptest.NewRequest(http.MethodPost, "/api/users", strings.NewReader(body))
		rec := httptest.NewRecorder()

		app.handleCreateUser(rec, req)

		if rec.Code != http.StatusCreated {
			t.Fatalf("code = %d, attendu %d", rec.Code, http.StatusCreated)
		}

		var got User
		if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
			t.Fatalf("reponse JSON invalide : %v", err)
		}
		if got.ID != 42 {
			t.Errorf("ID = %d, attendu 42", got.ID)
		}
		if got.Pseudo != "Thierry" {
			t.Errorf("Pseudo = %q, attendu Thierry", got.Pseudo)
		}
		if got.CreditBalance != welcomeCredits {
			t.Errorf("CreditBalance = %d, attendu %d", got.CreditBalance, welcomeCredits)
		}
	})

	t.Run("pseudo vide -> 400", func(t *testing.T) {
		app := &api{users: &fakeUserUseCase{
			registerFunc: func(ctx context.Context, pseudo, bio, city string) (User, error) {
				return User{}, ErrPseudoRequired
			},
		}}
		req := httptest.NewRequest(http.MethodPost, "/api/users", strings.NewReader(`{"pseudo":""}`))
		rec := httptest.NewRecorder()

		app.handleCreateUser(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("code = %d, attendu %d", rec.Code, http.StatusBadRequest)
		}
	})

	t.Run("JSON invalide -> 400", func(t *testing.T) {
		app := &api{users: &fakeUserUseCase{}}
		req := httptest.NewRequest(http.MethodPost, "/api/users", strings.NewReader(`{pas du json`))
		rec := httptest.NewRecorder()

		app.handleCreateUser(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("code = %d, attendu %d", rec.Code, http.StatusBadRequest)
		}
	})
}
