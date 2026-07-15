package main

import (
	"errors"
	"testing"
)

func TestNewService(t *testing.T) {
	cases := []struct {
		name       string
		title      string
		category   string
		city       string
		duration   int
		credits    int
		wantErr    error
		wantTitle  string
		wantActive bool
	}{
		{
			name:       "valid normalized ad",
			title:      "  Cours de Go  ",
			category:   "Informatique",
			city:       "  Paris  ",
			duration:   60,
			credits:    2,
			wantErr:    nil,
			wantTitle:  "Cours de Go",
			wantActive: true,
		},
		{
			name:     "empty title",
			title:    "   ",
			category: "Informatique",
			duration: 60,
			credits:  2,
			wantErr:  ErrServiceTitleRequired,
		},
		{
			name:     "category outside the closed list",
			title:    "Cours de Go",
			category: "Astrophysique",
			duration: 60,
			credits:  2,
			wantErr:  ErrServiceCategoryInvalid,
		},
		{
			name:     "zero duration",
			title:    "Cours de Go",
			category: "Informatique",
			duration: 0,
			credits:  2,
			wantErr:  ErrServiceDurationInvalid,
		},
		{
			name:     "negative credits",
			title:    "Cours de Go",
			category: "Informatique",
			duration: 60,
			credits:  -1,
			wantErr:  ErrServiceCreditsInvalid,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			service, err := NewService(5, testCase.title, "desc", testCase.category, testCase.city, testCase.duration, testCase.credits, true)

			if testCase.wantErr != nil {
				if !errors.Is(err, testCase.wantErr) {
					t.Fatalf("error = %v, want %v", err, testCase.wantErr)
				}

				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if service.Title != testCase.wantTitle {
				t.Errorf("Title = %q, want %q (trim)", service.Title, testCase.wantTitle)
			}
			if service.City != "Paris" {
				t.Errorf("City = %q, want Paris (trim)", service.City)
			}
			if service.ProviderID != 5 {
				t.Errorf("ProviderID = %d, want 5", service.ProviderID)
			}
			if service.Active != testCase.wantActive {
				t.Errorf("Active = %v, want %v", service.Active, testCase.wantActive)
			}
			if service.CreatedAt == "" {
				t.Error("CreatedAt must be set")
			}
		})
	}
}
