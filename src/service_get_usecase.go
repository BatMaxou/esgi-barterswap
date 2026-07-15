package main

import "context"

func (useCase *ServiceUseCase) Get(ctx context.Context, id int) (Service, error) {
	return useCase.services.FindByID(ctx, useCase.db.Executor(), id)
}
