package main

import "context"

func (useCase *ServiceUseCase) List(ctx context.Context, filter ServiceFilter) ([]Service, error) {
	return useCase.services.List(ctx, useCase.db.Executor(), filter)
}
