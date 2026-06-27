package main

import (
	"context"
	"net/http"
)

type userUseCase interface {
	Register(ctx context.Context, pseudo, bio, ville string) (User, error)
	GetProfile(ctx context.Context, id int) (User, error)
	Authenticate(ctx context.Context, id int) (User, error)
	UpdateProfile(ctx context.Context, actorID, targetID int, pseudo, bio, ville string) (User, error)
}

type api struct {
	users userUseCase
}

func (a *api) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /", handleIndex)

	mux.HandleFunc("POST /api/users", a.handleCreateUser)
	mux.HandleFunc("GET /api/users/{id}", a.handleGetUser)
	mux.HandleFunc("PUT /api/users/{id}", a.requireAuth(a.handleUpdateUser))
}
