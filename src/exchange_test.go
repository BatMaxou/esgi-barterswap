package main

import (
	"errors"
	"testing"
)

func TestNewExchange(t *testing.T) {
	testCases := []struct {
		name        string
		serviceID   int
		requesterID int
		ownerID     int
		wantErr     error
	}{
		{
			name:        "valid ids -> pending exchange with timestamps",
			serviceID:   3,
			requesterID: 2,
			ownerID:     1,
		},
		{
			name:        "zero service id -> ErrExchangeServiceIDInvalid",
			serviceID:   0,
			requesterID: 2,
			ownerID:     1,
			wantErr:     ErrExchangeServiceIDInvalid,
		},
		{
			name:        "negative service id -> ErrExchangeServiceIDInvalid",
			serviceID:   -1,
			requesterID: 2,
			ownerID:     1,
			wantErr:     ErrExchangeServiceIDInvalid,
		},
		{
			name:        "zero requester id -> ErrExchangeRequesterIDInvalid",
			serviceID:   3,
			requesterID: 0,
			ownerID:     1,
			wantErr:     ErrExchangeRequesterIDInvalid,
		},
		{
			name:        "zero owner id -> ErrExchangeOwnerIDInvalid",
			serviceID:   3,
			requesterID: 2,
			ownerID:     0,
			wantErr:     ErrExchangeOwnerIDInvalid,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			exchange, err := NewExchange(testCase.serviceID, testCase.requesterID, testCase.ownerID)
			if testCase.wantErr != nil {
				if !errors.Is(err, testCase.wantErr) {
					t.Fatalf("error = %v, want %v", err, testCase.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if exchange.ServiceID != testCase.serviceID {
				t.Errorf("ServiceID = %d, want %d", exchange.ServiceID, testCase.serviceID)
			}
			if exchange.RequesterID != testCase.requesterID {
				t.Errorf("RequesterID = %d, want %d", exchange.RequesterID, testCase.requesterID)
			}
			if exchange.OwnerID != testCase.ownerID {
				t.Errorf("OwnerID = %d, want %d", exchange.OwnerID, testCase.ownerID)
			}
			if exchange.Status != ExchangeStatusPending {
				t.Errorf("Status = %q, want %q", exchange.Status, ExchangeStatusPending)
			}
			if exchange.CreatedAt == "" {
				t.Error("CreatedAt must not be empty")
			}
			if exchange.UpdatedAt == "" {
				t.Error("UpdatedAt must not be empty")
			}
			if exchange.CreatedAt != exchange.UpdatedAt {
				t.Errorf("CreatedAt = %q, UpdatedAt = %q, want equal on creation", exchange.CreatedAt, exchange.UpdatedAt)
			}
		})
	}
}

func TestIsValidExchangeStatus(t *testing.T) {
	validStatuses := []string{
		ExchangeStatusPending,
		ExchangeStatusAccepted,
		ExchangeStatusRejected,
		ExchangeStatusCancelled,
		ExchangeStatusCompleted,
	}
	for _, status := range validStatuses {
		if !IsValidExchangeStatus(status) {
			t.Errorf("IsValidExchangeStatus(%q) = false, want true", status)
		}
	}

	invalidStatuses := []string{"", "unknown", "PENDING", "en cours"}
	for _, status := range invalidStatuses {
		if IsValidExchangeStatus(status) {
			t.Errorf("IsValidExchangeStatus(%q) = true, want false", status)
		}
	}
}
