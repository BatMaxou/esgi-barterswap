package main

import (
	"context"
	"errors"
	"testing"
)

func TestServiceUseCaseCreate(t *testing.T) {
	t.Run("valid creation sets the provider and activates the ad", func(t *testing.T) {
		services := &fakeServiceRepository{}
		useCase := NewServiceUseCase(&fakeDatabase{}, services, &fakeServiceExchangeRepository{})

		created, err := useCase.Create(context.Background(), 5, "  Cours de Go  ", "desc", "Informatique", "Paris", 60, 2)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !services.createCalled {
			t.Error("Create must be called")
		}
		if created.ProviderID != 5 {
			t.Errorf("ProviderID = %d, want 5", created.ProviderID)
		}
		if created.Title != "Cours de Go" {
			t.Errorf("Title = %q, want trimmed", created.Title)
		}
		if !created.Active {
			t.Error("a created ad must be active by default")
		}
		if created.ID != 42 {
			t.Errorf("ID = %d, want 42", created.ID)
		}
	})

	t.Run("invalid category -> ErrServiceCategoryInvalid without any write", func(t *testing.T) {
		services := &fakeServiceRepository{}
		useCase := NewServiceUseCase(&fakeDatabase{}, services, &fakeServiceExchangeRepository{})

		_, err := useCase.Create(context.Background(), 5, "Cours", "", "Truc", "", 60, 2)
		if !errors.Is(err, ErrServiceCategoryInvalid) {
			t.Fatalf("error = %v, want ErrServiceCategoryInvalid", err)
		}
		if services.createCalled {
			t.Error("no write must happen when validation fails")
		}
	})
}
