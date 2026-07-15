package main

import "context"

func (useCase *UserUseCase) Authenticate(ctx context.Context, id int) (User, error) {
	exec := useCase.db.Executor()

	return useCase.users.FindByID(ctx, exec, id)
}
