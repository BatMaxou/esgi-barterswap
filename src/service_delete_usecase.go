package main

import "context"

func (useCase *ServiceUseCase) Delete(ctx context.Context, actorID, serviceID int) error {
	exec := useCase.db.Executor()

	existing, err := useCase.services.FindByID(ctx, exec, serviceID)
	if err != nil {
		return err
	}
	if existing.ProviderID != actorID {
		return ErrForbidden
	}

	return useCase.services.Delete(ctx, exec, serviceID)
}
