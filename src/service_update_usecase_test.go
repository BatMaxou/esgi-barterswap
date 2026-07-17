package main

import (
	"context"
	"errors"
	"testing"
)

func TestServiceUseCaseUpdate(t *testing.T) {
	t.Run("update by the owner preserves id/createdAt and the active status", func(t *testing.T) {
		services := &fakeServiceRepository{service: Service{
			ID: 7, ProviderID: 5, Title: "Ancien", Category: "Informatique",
			DurationMinutes: 30, Credits: 1, Active: true, CreatedAt: "2026-01-01T00:00:00Z",
		}}
		useCase := NewServiceUseCase(&fakeDatabase{}, services, &fakeServiceExchangeRepository{})

		updated, err := useCase.Update(context.Background(), 5, 7, "Nouveau", "desc", "Cuisine", "Lyon", 90, 3, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !services.updateCalled {
			t.Error("Update must be called")
		}
		if updated.Title != "Nouveau" || updated.Category != "Cuisine" || updated.Credits != 3 {
			t.Errorf("fields not updated: %+v", updated)
		}
		if updated.ID != 7 {
			t.Errorf("ID = %d, must be preserved (7)", updated.ID)
		}
		if updated.CreatedAt != "2026-01-01T00:00:00Z" {
			t.Errorf("CreatedAt = %q, must be preserved", updated.CreatedAt)
		}
		if !updated.Active {
			t.Error("Active must stay unchanged (true) when not provided")
		}
	})

	t.Run("deactivation via the active field", func(t *testing.T) {
		services := &fakeServiceRepository{service: Service{
			ID: 7, ProviderID: 5, Title: "Ancien", Category: "Informatique",
			DurationMinutes: 30, Credits: 1, Active: true, CreatedAt: "2026-01-01T00:00:00Z",
		}}
		useCase := NewServiceUseCase(&fakeDatabase{}, services, &fakeServiceExchangeRepository{})

		inactive := false
		updated, err := useCase.Update(context.Background(), 5, 7, "Ancien", "", "Informatique", "", 30, 1, &inactive)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if updated.Active {
			t.Error("Active must become false")
		}
	})

	t.Run("another user's ad -> ErrForbidden", func(t *testing.T) {
		services := &fakeServiceRepository{service: Service{ID: 7, ProviderID: 5}}
		useCase := NewServiceUseCase(&fakeDatabase{}, services, &fakeServiceExchangeRepository{})

		_, err := useCase.Update(context.Background(), 9, 7, "X", "", "Informatique", "", 30, 1, nil)
		if !errors.Is(err, ErrForbidden) {
			t.Fatalf("error = %v, want ErrForbidden", err)
		}
		if services.updateCalled {
			t.Error("no write must happen on a forbidden access")
		}
	})

	t.Run("ad not found -> ErrServiceNotFound", func(t *testing.T) {
		services := &fakeServiceRepository{findErr: ErrServiceNotFound}
		useCase := NewServiceUseCase(&fakeDatabase{}, services, &fakeServiceExchangeRepository{})

		_, err := useCase.Update(context.Background(), 5, 999, "X", "", "Informatique", "", 30, 1, nil)
		if !errors.Is(err, ErrServiceNotFound) {
			t.Fatalf("error = %v, want ErrServiceNotFound", err)
		}
	})
}
