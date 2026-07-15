package main

import (
	"context"
	"testing"
)

func TestServiceUseCaseList(t *testing.T) {
	t.Run("filters are forwarded to the repository", func(t *testing.T) {
		services := &fakeServiceRepository{services: []Service{{ID: 1, Title: "Cours de Go"}}}
		useCase := NewServiceUseCase(&fakeDatabase{}, services)

		filter := ServiceFilter{Category: "Informatique", City: "Paris", Search: "Go"}
		got, err := useCase.List(context.Background(), filter)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got) != 1 {
			t.Fatalf("services = %d, want 1", len(got))
		}
		if services.filter != filter {
			t.Errorf("forwarded filter = %+v, want %+v", services.filter, filter)
		}
	})
}
