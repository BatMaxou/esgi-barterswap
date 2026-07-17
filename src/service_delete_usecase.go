package main

import "context"

func (useCase *ServiceUseCase) Delete(ctx context.Context, actorID, serviceID int) error {
	existing, err := useCase.services.FindByID(ctx, useCase.db.Executor(), serviceID)
	if err != nil {
		return err
	}
	if existing.ProviderID != actorID {
		return ErrForbidden
	}

	// The check and the delete share a transaction: an exchange created in
	// between would otherwise break the exchanges -> services foreign key.
	return useCase.db.WithinTransaction(ctx, func(exec dbExecutor) error {
		referenced, err := useCase.exchanges.HasAnyForService(ctx, exec, serviceID)
		if err != nil {
			return err
		}
		if referenced {
			return ErrServiceHasExchanges
		}

		return useCase.services.Delete(ctx, exec, serviceID)
	})
}
