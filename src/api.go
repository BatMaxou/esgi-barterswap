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

type exchangeUseCase interface {
	Create(ctx context.Context, requesterID, serviceID int) (Exchange, error)
	List(ctx context.Context, actorID int, status string) ([]Exchange, error)
	Get(ctx context.Context, actorID, exchangeID int) (Exchange, error)
	Accept(ctx context.Context, actorID, exchangeID int) (Exchange, error)
	Reject(ctx context.Context, actorID, exchangeID int) (Exchange, error)
	Complete(ctx context.Context, actorID, exchangeID int) (Exchange, error)
	Cancel(ctx context.Context, actorID, exchangeID int) (Exchange, error)
}

type reviewUseCase interface {
	Create(ctx context.Context, actorID, exchangeID, rating int, comment string) (Review, error)
	ListForUser(ctx context.Context, userID int) ([]Review, error)
	ListForService(ctx context.Context, serviceID int) ([]Review, error)
}

type userStatsUseCase interface {
	Get(ctx context.Context, userID int) (UserStats, error)
}

type api struct {
	users     userUseCase
	skills    skillUseCase
	services  serviceUseCase
	exchanges exchangeUseCase
	reviews   reviewUseCase
	stats     userStatsUseCase
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

	mux.HandleFunc("POST /api/exchanges", a.requireAuth(a.handleCreateExchange))
	mux.HandleFunc("GET /api/exchanges", a.requireAuth(a.handleListExchanges))
	mux.HandleFunc("GET /api/exchanges/{id}", a.requireAuth(a.handleGetExchange))
	mux.HandleFunc("PUT /api/exchanges/{id}/accept", a.requireAuth(a.handleAcceptExchange))
	mux.HandleFunc("PUT /api/exchanges/{id}/reject", a.requireAuth(a.handleRejectExchange))
	mux.HandleFunc("PUT /api/exchanges/{id}/complete", a.requireAuth(a.handleCompleteExchange))
	mux.HandleFunc("PUT /api/exchanges/{id}/cancel", a.requireAuth(a.handleCancelExchange))

	mux.HandleFunc("POST /api/exchanges/{id}/review", a.requireAuth(a.handleCreateReview))
	mux.HandleFunc("GET /api/users/{id}/reviews", a.handleListUserReviews)
	mux.HandleFunc("GET /api/users/{id}/stats", a.handleGetUserStats)
	mux.HandleFunc("GET /api/services/{id}/reviews", a.handleListServiceReviews)
}
