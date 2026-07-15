package main

import (
	"errors"
	"testing"
)

func TestNewUser(t *testing.T) {
	tests := []struct {
		name       string
		pseudo     string
		bio        string
		city       string
		wantErr    error
		wantPseudo string
	}{
		{name: "valid pseudo", pseudo: "Thierry", bio: "my bio", city: "Paris", wantPseudo: "Thierry"},
		{name: "pseudo surrounded by spaces", pseudo: "  Alice  ", wantPseudo: "Alice"},
		{name: "empty pseudo", pseudo: "", wantErr: ErrPseudoRequired},
		{name: "whitespace-only pseudo", pseudo: "   ", wantErr: ErrPseudoRequired},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			user, err := NewUser(testCase.pseudo, testCase.bio, testCase.city)

			if testCase.wantErr != nil {
				if !errors.Is(err, testCase.wantErr) {
					t.Fatalf("error = %v, want %v", err, testCase.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if user.Pseudo != testCase.wantPseudo {
				t.Errorf("Pseudo = %q, want %q", user.Pseudo, testCase.wantPseudo)
			}
			if user.CreditBalance != 0 {
				t.Errorf("CreditBalance = %d, want 0 (welcome credits are granted by the use case)", user.CreditBalance)
			}
			if user.Skills == nil {
				t.Error("Skills must be a non-nil empty slice")
			}
			if user.CreatedAt == "" {
				t.Error("CreatedAt must not be empty")
			}
		})
	}
}
