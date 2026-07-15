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

	t.Run("setting your own skills -> 200", func(t *testing.T) {
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
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
		}
		var got []Skill
		if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
			t.Fatalf("invalid JSON response: %v", err)
		}
		if len(got) != 1 || got[0].Level != "expert" {
			t.Errorf("skills = %+v, want [expert]", got)
		}
	})

	t.Run("another user's skills -> 403", func(t *testing.T) {
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
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusForbidden)
		}
	})

	t.Run("invalid level -> 400", func(t *testing.T) {
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
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
	})

	t.Run("without authenticated user -> 401", func(t *testing.T) {
		app := &api{skills: &fakeSkillUseCase{}}

		req := httptest.NewRequest(http.MethodPut, "/api/users/5/skills", strings.NewReader(`[]`))
		req.SetPathValue("id", "5")
		rec := httptest.NewRecorder()

		app.handleDefineUserSkills(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
		}
	})

	t.Run("invalid JSON -> 400", func(t *testing.T) {
		app := &api{skills: &fakeSkillUseCase{}}

		req := httptest.NewRequest(http.MethodPut, "/api/users/5/skills", strings.NewReader(`{not json`))
		req.SetPathValue("id", "5")
		req = withCurrentUser(req, User{ID: 5})
		rec := httptest.NewRecorder()

		app.handleDefineUserSkills(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
	})
}
