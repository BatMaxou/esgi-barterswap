package main

import "context"

type fakeExchangeRepository struct {
	exchange           Exchange
	exchanges          []Exchange
	createCalled       bool
	updateStatusCalled bool
	hasActiveCalled    bool
	hasActive          bool
	created            Exchange
	updatedStatus      string
	updatedAt          string
	filter             ExchangeFilter
	findErr            error
	createErr          error
	updateStatusErr    error
	hasActiveErr       error
	listErr            error
	serviceExchanges   []Exchange
	listByServiceErr   error
	completedCount     int
}

func (fake *fakeExchangeRepository) Create(ctx context.Context, exec dbExecutor, exchange Exchange) (Exchange, error) {
	fake.createCalled = true
	if fake.createErr != nil {
		return Exchange{}, fake.createErr
	}
	exchange.ID = 9
	fake.created = exchange
	return exchange, nil
}

func (fake *fakeExchangeRepository) FindByID(ctx context.Context, exec dbExecutor, id int) (Exchange, error) {
	if fake.findErr != nil {
		return Exchange{}, fake.findErr
	}
	return fake.exchange, nil
}

func (fake *fakeExchangeRepository) UpdateStatus(ctx context.Context, exec dbExecutor, id int, status, updatedAt string) (Exchange, error) {
	fake.updateStatusCalled = true
	fake.updatedStatus = status
	fake.updatedAt = updatedAt
	if fake.updateStatusErr != nil {
		return Exchange{}, fake.updateStatusErr
	}
	fake.exchange.Status = status
	fake.exchange.UpdatedAt = updatedAt
	return fake.exchange, nil
}

func (fake *fakeExchangeRepository) HasActiveForService(ctx context.Context, exec dbExecutor, serviceID int) (bool, error) {
	fake.hasActiveCalled = true
	if fake.hasActiveErr != nil {
		return false, fake.hasActiveErr
	}
	return fake.hasActive, nil
}

func (fake *fakeExchangeRepository) List(ctx context.Context, exec dbExecutor, filter ExchangeFilter) ([]Exchange, error) {
	fake.filter = filter
	if fake.listErr != nil {
		return nil, fake.listErr
	}
	return fake.exchanges, nil
}

func (fake *fakeExchangeRepository) ListByServiceID(ctx context.Context, exec dbExecutor, serviceID int) ([]Exchange, error) {
	if fake.listByServiceErr != nil {
		return nil, fake.listByServiceErr
	}
	return fake.serviceExchanges, nil
}

func (fake *fakeExchangeRepository) CountCompletedByUserID(ctx context.Context, exec dbExecutor, userID int) (int, error) {
	return fake.completedCount, nil
}
