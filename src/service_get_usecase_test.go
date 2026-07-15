package main

import (
	"context"
	"errors"
	"testing"
)

func TestServiceUseCaseGet(t *testing.T) {
	t.Run("service not found", func(t *testing.T) {
		services := &fakeServiceRepository{findErr: ErrServiceNotFound}
		useCase := NewServiceUseCase(&fakeDatabase{}, services)

		_, err := useCase.Get(context.Background(), 999)
		if !errors.Is(err, ErrServiceNotFound) {
			t.Fatalf("error = %v, want ErrServiceNotFound", err)
		}
	})
}
