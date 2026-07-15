package main

import (
	"errors"
	"time"
)

const (
	ExchangeStatusPending   = "pending"
	ExchangeStatusAccepted  = "accepted"
	ExchangeStatusRejected  = "rejected"
	ExchangeStatusCancelled = "cancelled"
	ExchangeStatusCompleted = "completed"
)

var ErrExchangeNotFound = errors.New("exchange not found")

var ErrExchangeServiceIDInvalid = errors.New("service id must be strictly positive")

var ErrExchangeRequesterIDInvalid = errors.New("requester id must be strictly positive")

var ErrExchangeOwnerIDInvalid = errors.New("owner id must be strictly positive")

var ErrExchangeSelfRequest = errors.New("cannot request an exchange on your own service")

var ErrExchangeInsufficientCredits = errors.New("insufficient credits for this exchange")

var ErrExchangeServiceUnavailable = errors.New("service already has an active exchange")

var ErrExchangeServiceInactive = errors.New("service is not active")

var ErrExchangeInvalidTransition = errors.New("invalid exchange status transition")

var ErrExchangeStatusInvalid = errors.New("invalid exchange status filter")

type Exchange struct {
	ID          int    `json:"id"`
	ServiceID   int    `json:"service_id"`
	RequesterID int    `json:"requester_id"`
	OwnerID     int    `json:"owner_id"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

var validExchangeStatuses = map[string]bool{
	ExchangeStatusPending:   true,
	ExchangeStatusAccepted:  true,
	ExchangeStatusRejected:  true,
	ExchangeStatusCancelled: true,
	ExchangeStatusCompleted: true,
}

// NewExchange validates the identifiers and returns an Exchange ready to be
// persisted with status pending. Business rules such as self-requests or credit
// checks belong in the use case.
func NewExchange(serviceID, requesterID, ownerID int) (Exchange, error) {
	if serviceID <= 0 {
		return Exchange{}, ErrExchangeServiceIDInvalid
	}
	if requesterID <= 0 {
		return Exchange{}, ErrExchangeRequesterIDInvalid
	}
	if ownerID <= 0 {
		return Exchange{}, ErrExchangeOwnerIDInvalid
	}

	now := time.Now().UTC().Format(time.RFC3339)

	return Exchange{
		ServiceID:   serviceID,
		RequesterID: requesterID,
		OwnerID:     ownerID,
		Status:      ExchangeStatusPending,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

func IsValidExchangeStatus(status string) bool {
	return validExchangeStatuses[status]
}
