package main

import (
	"context"
	"net/http"
)

type userUseCase interface {
	Register(ctx context.Context, pseudo, bio, city string) (User, error)
	GetProfile(ctx context.Context, id int) (User, error)
	Authenticate(ctx context.Context, id int) (User, error)
	UpdateProfile(ctx context.Context, actorID, targetID int, pseudo, bio, city string) (User, error)
}

type skillUseCase interface {
	ListSkills(ctx context.Context, userID int) ([]Skill, error)
	DefineSkills(ctx context.Context, actorID, targetID int, skills []Skill) ([]Skill, error)
}

type serviceUseCase interface {
	Create(ctx context.Context, providerID int, title, description, category, city string, durationMinutes, credits int) (Service, error)
	Get(ctx context.Context, id int) (Service, error)
	List(ctx context.Context, filter ServiceFilter) ([]Service, error)
	Update(ctx context.Context, actorID, serviceID int, title, description, category, city string, durationMinutes, credits int, active *bool) (Service, error)
	Delete(ctx context.Context, actorID, serviceID int) error
}

type api struct {
	users    userUseCase
	skills   skillUseCase
	services serviceUseCase
}

func (a *api) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /", handleIndex)

	mux.HandleFunc("POST /api/users", a.handleCreateUser)
	mux.HandleFunc("GET /api/users/{id}", a.handleGetUser)
	mux.HandleFunc("PUT /api/users/{id}", a.requireAuth(a.handleUpdateUser))
	mux.HandleFunc("GET /api/users/{id}/skills", a.handleGetUserSkills)
	mux.HandleFunc("PUT /api/users/{id}/skills", a.requireAuth(a.handleDefineUserSkills))

	mux.HandleFunc("GET /api/services", a.handleListServices)
	mux.HandleFunc("POST /api/services", a.requireAuth(a.handleCreateService))
	mux.HandleFunc("GET /api/services/{id}", a.handleGetService)
	mux.HandleFunc("PUT /api/services/{id}", a.requireAuth(a.handleUpdateService))
	mux.HandleFunc("DELETE /api/services/{id}", a.requireAuth(a.handleDeleteService))
}
