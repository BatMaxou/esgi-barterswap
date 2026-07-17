package main

import (
	"errors"
	"strings"
	"time"
)

var ErrServiceTitleRequired = errors.New("service title is required")

var ErrServiceCategoryInvalid = errors.New("invalid category (not in the category list)")

var ErrServiceDurationInvalid = errors.New("duration in minutes must be strictly positive")

var ErrServiceCreditsInvalid = errors.New("credit cost must be strictly positive")

var ErrServiceNotFound = errors.New("service not found")

var ErrServiceHasExchanges = errors.New("service is referenced by an exchange and cannot be deleted")

type Service struct {
	ID              int    `json:"id"`
	ProviderID      int    `json:"provider_id"`
	Title           string `json:"title"`
	Description     string `json:"description,omitempty"`
	Category        string `json:"category"`
	DurationMinutes int    `json:"duration_minutes"`
	Credits         int    `json:"credits"`
	City            string `json:"city,omitempty"`
	Active          bool   `json:"active"`
	CreatedAt       string `json:"created_at"`
}

// NewService normalizes and validates the fields of an ad, then returns a
// Service ready to be persisted. Categories are restricted to the closed list
// required by the subject (shared with skills).
func NewService(providerID int, title, description, category, city string, durationMinutes, credits int, active bool) (Service, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return Service{}, ErrServiceTitleRequired
	}

	category = strings.TrimSpace(category)
	if !validCategories[category] {
		return Service{}, ErrServiceCategoryInvalid
	}

	if durationMinutes <= 0 {
		return Service{}, ErrServiceDurationInvalid
	}

	if credits <= 0 {
		return Service{}, ErrServiceCreditsInvalid
	}

	return Service{
		ProviderID:      providerID,
		Title:           title,
		Description:     strings.TrimSpace(description),
		Category:        category,
		DurationMinutes: durationMinutes,
		Credits:         credits,
		City:            strings.TrimSpace(city),
		Active:          active,
		CreatedAt:       time.Now().UTC().Format(time.RFC3339),
	}, nil
}
