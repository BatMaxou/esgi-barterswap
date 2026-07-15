package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type fakeUserUseCase struct {
	registerFunc      func(ctx context.Context, pseudo, bio, city string) (User, error)
	getProfileFunc    func(ctx context.Context, id int) (User, error)
	authenticateFunc  func(ctx context.Context, id int) (User, error)
	updateProfileFunc func(ctx context.Context, actorID, targetID int, pseudo, bio, city string) (User, error)
}

func (fake *fakeUserUseCase) Register(ctx context.Context, pseudo, bio, city string) (User, error) {
	return fake.registerFunc(ctx, pseudo, bio, city)
}

func (fake *fakeUserUseCase) GetProfile(ctx context.Context, id int) (User, error) {
	return fake.getProfileFunc(ctx, id)
}

func (fake *fakeUserUseCase) Authenticate(ctx context.Context, id int) (User, error) {
	return fake.authenticateFunc(ctx, id)
}

func (fake *fakeUserUseCase) UpdateProfile(ctx context.Context, actorID, targetID int, pseudo, bio, city string) (User, error) {
	return fake.updateProfileFunc(ctx, actorID, targetID, pseudo, bio, city)
}

type fakeSkillUseCase struct {
	listSkillsFunc   func(ctx context.Context, userID int) ([]Skill, error)
	defineSkillsFunc func(ctx context.Context, actorID, targetID int, skills []Skill) ([]Skill, error)
}

func (fake *fakeSkillUseCase) ListSkills(ctx context.Context, userID int) ([]Skill, error) {
	return fake.listSkillsFunc(ctx, userID)
}

func (fake *fakeSkillUseCase) DefineSkills(ctx context.Context, actorID, targetID int, skills []Skill) ([]Skill, error) {
	return fake.defineSkillsFunc(ctx, actorID, targetID, skills)
}

func TestUpdateUserRouting(t *testing.T) {
	app := &api{users: &fakeUserUseCase{
		authenticateFunc: func(ctx context.Context, id int) (User, error) {
			return User{ID: id, Pseudo: "Thierry"}, nil
		},
		updateProfileFunc: func(ctx context.Context, actorID, targetID int, pseudo, bio, city string) (User, error) {
			return User{ID: targetID, Pseudo: pseudo}, nil
		},
	}}
	mux := http.NewServeMux()
	app.registerRoutes(mux)

	t.Run("without X-User-ID header -> 401", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/api/users/5", strings.NewReader(`{"pseudo":"Thierry"}`))
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
		}
	})

	t.Run("with valid X-User-ID header -> 200", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/api/users/5", strings.NewReader(`{"pseudo":"Thierry"}`))
		req.Header.Set("X-User-ID", "5")
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
		}
	})
}
