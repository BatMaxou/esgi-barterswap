package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandleDefineUserSkills(t *testing.T) {
	withCurrentUser := func(req *http.Request, user User) *http.Request {
		ctx := context.WithValue(req.Context(), currentUserKey, user)
		return req.WithContext(ctx)
	}

	t.Run("definition de ses competences -> 200", func(t *testing.T) {
		app := &api{skills: &fakeSkillUseCase{
			defineSkillsFunc: func(ctx context.Context, actorID, targetID int, skills []Skill) ([]Skill, error) {
				return skills, nil
			},
		}}

		body := `[{"name":"Jardinage","level":"expert"}]`
		req := httptest.NewRequest(http.MethodPut, "/api/users/5/skills", strings.NewReader(body))
		req.SetPathValue("id", "5")
		req = withCurrentUser(req, User{ID: 5})
		rec := httptest.NewRecorder()

		app.handleDefineUserSkills(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("code = %d, attendu %d", rec.Code, http.StatusOK)
		}
		var got []Skill
		if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
			t.Fatalf("reponse JSON invalide : %v", err)
		}
		if len(got) != 1 || got[0].Level != "expert" {
			t.Errorf("competences = %+v, attendu [expert]", got)
		}
	})

	t.Run("competences d'un autre utilisateur -> 403", func(t *testing.T) {
		app := &api{skills: &fakeSkillUseCase{
			defineSkillsFunc: func(ctx context.Context, actorID, targetID int, skills []Skill) ([]Skill, error) {
				return nil, ErrForbidden
			},
		}}

		req := httptest.NewRequest(http.MethodPut, "/api/users/9/skills", strings.NewReader(`[{"name":"Jardinage","level":"expert"}]`))
		req.SetPathValue("id", "9")
		req = withCurrentUser(req, User{ID: 5})
		rec := httptest.NewRecorder()

		app.handleDefineUserSkills(rec, req)

		if rec.Code != http.StatusForbidden {
			t.Fatalf("code = %d, attendu %d", rec.Code, http.StatusForbidden)
		}
	})

	t.Run("level invalide -> 400", func(t *testing.T) {
		app := &api{skills: &fakeSkillUseCase{
			defineSkillsFunc: func(ctx context.Context, actorID, targetID int, skills []Skill) ([]Skill, error) {
				return nil, ErrSkillLevelInvalid
			},
		}}

		req := httptest.NewRequest(http.MethodPut, "/api/users/5/skills", strings.NewReader(`[{"name":"Jardinage","level":"maitre"}]`))
		req.SetPathValue("id", "5")
		req = withCurrentUser(req, User{ID: 5})
		rec := httptest.NewRecorder()

		app.handleDefineUserSkills(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("code = %d, attendu %d", rec.Code, http.StatusBadRequest)
		}
	})

	t.Run("sans utilisateur authentifie -> 401", func(t *testing.T) {
		app := &api{skills: &fakeSkillUseCase{}}

		req := httptest.NewRequest(http.MethodPut, "/api/users/5/skills", strings.NewReader(`[]`))
		req.SetPathValue("id", "5")
		rec := httptest.NewRecorder()

		app.handleDefineUserSkills(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("code = %d, attendu %d", rec.Code, http.StatusUnauthorized)
		}
	})

	t.Run("JSON invalide -> 400", func(t *testing.T) {
		app := &api{skills: &fakeSkillUseCase{}}

		req := httptest.NewRequest(http.MethodPut, "/api/users/5/skills", strings.NewReader(`{pas du json`))
		req.SetPathValue("id", "5")
		req = withCurrentUser(req, User{ID: 5})
		rec := httptest.NewRecorder()

		app.handleDefineUserSkills(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("code = %d, attendu %d", rec.Code, http.StatusBadRequest)
		}
	})
}
