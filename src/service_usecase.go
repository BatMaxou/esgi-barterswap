package main

import "context"

type serviceRepository interface {
	Create(ctx context.Context, exec dbExecutor, service Service) (Service, error)
	Update(ctx context.Context, exec dbExecutor, service Service) (Service, error)
	Delete(ctx context.Context, exec dbExecutor, id int) error
	FindByID(ctx context.Context, exec dbExecutor, id int) (Service, error)
	List(ctx context.Context, exec dbExecutor, filter ServiceFilter) ([]Service, error)
}

type serviceExchangeRepository interface {
	HasAnyForService(ctx context.Context, exec dbExecutor, serviceID int) (bool, error)
}

// ServiceFilter carries the search criteria applied server-side.
// An empty criterion is ignored.
type ServiceFilter struct {
	Category string
	City     string
	Search   string
}

type ServiceUseCase struct {
	db        database
	services  serviceRepository
	exchanges serviceExchangeRepository
}

func NewServiceUseCase(db database, services serviceRepository, exchanges serviceExchangeRepository) *ServiceUseCase {
	return &ServiceUseCase{
		db:        db,
		services:  services,
		exchanges: exchanges,
	}
}
