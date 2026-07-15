package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleGetUserSkills(t *testing.T) {
	t.Run("existing skills -> 200", func(t *testing.T) {
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
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
		}
		var got []Skill
		if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
			t.Fatalf("invalid JSON response: %v", err)
		}
		if len(got) != 1 || got[0].Name != "Jardinage" {
			t.Errorf("skills = %+v, want [Jardinage]", got)
		}
	})

	t.Run("user not found -> 404", func(t *testing.T) {
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
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
		}
	})

	t.Run("invalid identifier -> 400", func(t *testing.T) {
		app := &api{skills: &fakeSkillUseCase{}}

		req := httptest.NewRequest(http.MethodGet, "/api/users/abc/skills", nil)
		req.SetPathValue("id", "abc")
		rec := httptest.NewRecorder()

		app.handleGetUserSkills(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
	})
}
