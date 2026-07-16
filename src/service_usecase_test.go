package main

import "context"

// fakeServiceRepository is shared by the ServiceUseCase operation tests.
type fakeServiceRepository struct {
	service      Service
	services     []Service
	createCalled bool
	updateCalled bool
	deleteCalled bool
	created      Service
	updated      Service
	deletedID    int
	filter       ServiceFilter
	findErr      error
	createErr    error
	updateErr    error
	deleteErr    error
	listErr            error
	activeCount        int
}

func (fake *fakeServiceRepository) Create(ctx context.Context, exec dbExecutor, service Service) (Service, error) {
	fake.createCalled = true
	if fake.createErr != nil {
		return Service{}, fake.createErr
	}
	service.ID = 42
	fake.created = service
	return service, nil
}

func (fake *fakeServiceRepository) Update(ctx context.Context, exec dbExecutor, service Service) (Service, error) {
	fake.updateCalled = true
	if fake.updateErr != nil {
		return Service{}, fake.updateErr
	}
	fake.updated = service
	return service, nil
}

func (fake *fakeServiceRepository) Delete(ctx context.Context, exec dbExecutor, id int) error {
	fake.deleteCalled = true
	if fake.deleteErr != nil {
		return fake.deleteErr
	}
	fake.deletedID = id
	return nil
}

func (fake *fakeServiceRepository) FindByID(ctx context.Context, exec dbExecutor, id int) (Service, error) {
	if fake.findErr != nil {
		return Service{}, fake.findErr
	}
	return fake.service, nil
}

func (fake *fakeServiceRepository) List(ctx context.Context, exec dbExecutor, filter ServiceFilter) ([]Service, error) {
	fake.filter = filter
	if fake.listErr != nil {
		return nil, fake.listErr
	}
	return fake.services, nil
}

func (fake *fakeServiceRepository) CountActiveByProviderID(ctx context.Context, exec dbExecutor, providerID int) (int, error) {
	return fake.activeCount, nil
}
