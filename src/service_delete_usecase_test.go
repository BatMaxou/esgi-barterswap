package main

import (
	"context"
	"errors"
	"testing"
)

func TestServiceUseCaseDelete(t *testing.T) {
	t.Run("delete by the owner", func(t *testing.T) {
		services := &fakeServiceRepository{service: Service{ID: 7, ProviderID: 5}}
		useCase := NewServiceUseCase(&fakeDatabase{}, services)

		if err := useCase.Delete(context.Background(), 5, 7); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !services.deleteCalled || services.deletedID != 7 {
			t.Errorf("Delete must be called with id 7 (called=%v, id=%d)", services.deleteCalled, services.deletedID)
		}
	})

	t.Run("another user's ad -> ErrForbidden", func(t *testing.T) {
		services := &fakeServiceRepository{service: Service{ID: 7, ProviderID: 5}}
		useCase := NewServiceUseCase(&fakeDatabase{}, services)

		err := useCase.Delete(context.Background(), 9, 7)
		if !errors.Is(err, ErrForbidden) {
			t.Fatalf("error = %v, want ErrForbidden", err)
		}
		if services.deleteCalled {
			t.Error("no delete must happen on a forbidden access")
		}
	})

	t.Run("ad not found -> ErrServiceNotFound", func(t *testing.T) {
		services := &fakeServiceRepository{findErr: ErrServiceNotFound}
		useCase := NewServiceUseCase(&fakeDatabase{}, services)

		err := useCase.Delete(context.Background(), 5, 999)
		if !errors.Is(err, ErrServiceNotFound) {
			t.Fatalf("error = %v, want ErrServiceNotFound", err)
		}
	})
}
