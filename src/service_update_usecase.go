package main

import (
	"context"
	"fmt"
)

func (useCase *ServiceUseCase) Update(ctx context.Context, actorID, serviceID int, title, description, category, city string, durationMinutes, credits int, active *bool) (Service, error) {
	exec := useCase.db.Executor()

	existing, err := useCase.services.FindByID(ctx, exec, serviceID)
	if err != nil {
		return Service{}, err
	}
	if existing.ProviderID != actorID {
		return Service{}, ErrForbidden
	}

	nextActive := existing.Active
	if active != nil {
		nextActive = *active
	}

	candidate, err := NewService(existing.ProviderID, title, description, category, city, durationMinutes, credits, nextActive)
	if err != nil {
		return Service{}, err
	}
	candidate.ID = existing.ID
	candidate.CreatedAt = existing.CreatedAt

	updated, err := useCase.services.Update(ctx, exec, candidate)
	if err != nil {
		return Service{}, fmt.Errorf("update service: %w", err)
	}

	return updated, nil
}
