package main

import "context"

func (useCase *ReviewUseCase) ListForUser(ctx context.Context, userID int) ([]Review, error) {
	exec := useCase.db.Executor()

	if _, err := useCase.users.FindByID(ctx, exec, userID); err != nil {
		return nil, err
	}

	return useCase.reviews.ListByTargetUserID(ctx, exec, userID)
}
