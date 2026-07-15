package main

import (
	"context"
	"time"
)

type exchangeRepository interface {
	Create(ctx context.Context, exec dbExecutor, exchange Exchange) (Exchange, error)
	FindByID(ctx context.Context, exec dbExecutor, id int) (Exchange, error)
	UpdateStatus(ctx context.Context, exec dbExecutor, id int, status, updatedAt string) (Exchange, error)
	HasActiveForService(ctx context.Context, exec dbExecutor, serviceID int) (bool, error)
	List(ctx context.Context, exec dbExecutor, filter ExchangeFilter) ([]Exchange, error)
	ListByServiceID(ctx context.Context, exec dbExecutor, serviceID int) ([]Exchange, error)
}

type ExchangeUseCase struct {
	db                 database
	exchanges          exchangeRepository
	services           serviceRepository
	creditTransactions creditTransactionRepository
}

func NewExchangeUseCase(
	db database,
	exchanges exchangeRepository,
	services serviceRepository,
	creditTransactions creditTransactionRepository,
) *ExchangeUseCase {
	return &ExchangeUseCase{
		db:                 db,
		exchanges:          exchanges,
		services:           services,
		creditTransactions: creditTransactions,
	}
}

func (useCase *ExchangeUseCase) ensureActorCanView(actorID int, exchange Exchange) error {
	if actorID != exchange.RequesterID && actorID != exchange.OwnerID {
		return ErrForbidden
	}

	return nil
}

func creditTransactionForExchange(userID, exchangeID, amount int, transactionType string) CreditTransaction {
	return CreditTransaction{
		UserID:     userID,
		ExchangeID: exchangeID,
		Amount:     amount,
		Type:       transactionType,
		CreatedAt:  time.Now().UTC().Format(time.RFC3339),
	}
}

func exchangeUpdatedNow() string {
	return time.Now().UTC().Format(time.RFC3339)
}
