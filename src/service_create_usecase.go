package main

import (
	"context"
	"fmt"
)

func (useCase *ServiceUseCase) Create(ctx context.Context, providerID int, title, description, category, city string, durationMinutes, credits int) (Service, error) {
	service, err := NewService(providerID, title, description, category, city, durationMinutes, credits, true)
	if err != nil {
		return Service{}, err
	}

	created, err := useCase.services.Create(ctx, useCase.db.Executor(), service)
	if err != nil {
		return Service{}, fmt.Errorf("create service: %w", err)
	}

	return created, nil
}
