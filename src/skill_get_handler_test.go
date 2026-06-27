package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleGetUserSkills(t *testing.T) {
	t.Run("competences existantes -> 200", func(t *testing.T) {
		app := &api{skills: &fakeSkillUseCase{
			listSkillsFunc: func(ctx context.Context, userID int) ([]Skill, error) {
				return []Skill{{Name: "Jardinage", Level: "expert"}}, nil
			},
		}}

		req := httptest.NewRequest(http.MethodGet, "/api/users/5/skills", nil)
		req.SetPathValue("id", "5")
		rec := httptest.NewRecorder()

		app.handleGetUserSkills(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("code = %d, attendu %d", rec.Code, http.StatusOK)
		}
		var got []Skill
		if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
			t.Fatalf("reponse JSON invalide : %v", err)
		}
		if len(got) != 1 || got[0].Name != "Jardinage" {
			t.Errorf("competences = %+v, attendu [Jardinage]", got)
		}
	})

	t.Run("utilisateur introuvable -> 404", func(t *testing.T) {
		app := &api{skills: &fakeSkillUseCase{
			listSkillsFunc: func(ctx context.Context, userID int) ([]Skill, error) {
				return nil, ErrUserNotFound
			},
		}}

		req := httptest.NewRequest(http.MethodGet, "/api/users/999/skills", nil)
		req.SetPathValue("id", "999")
		rec := httptest.NewRecorder()

		app.handleGetUserSkills(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Fatalf("code = %d, attendu %d", rec.Code, http.StatusNotFound)
		}
	})

	t.Run("identifiant invalide -> 400", func(t *testing.T) {
		app := &api{skills: &fakeSkillUseCase{}}

		req := httptest.NewRequest(http.MethodGet, "/api/users/abc/skills", nil)
		req.SetPathValue("id", "abc")
		rec := httptest.NewRecorder()

		app.handleGetUserSkills(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("code = %d, attendu %d", rec.Code, http.StatusBadRequest)
		}
	})
}
