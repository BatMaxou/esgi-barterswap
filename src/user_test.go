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
		{name: "pseudo valide", pseudo: "Thierry", bio: "ma bio", city: "Paris", wantPseudo: "Thierry"},
		{name: "pseudo entoure d'espaces", pseudo: "  Alice  ", wantPseudo: "Alice"},
		{name: "pseudo vide", pseudo: "", wantErr: ErrPseudoRequired},
		{name: "pseudo espaces uniquement", pseudo: "   ", wantErr: ErrPseudoRequired},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			user, err := NewUser(testCase.pseudo, testCase.bio, testCase.city)

			if testCase.wantErr != nil {
				if !errors.Is(err, testCase.wantErr) {
					t.Fatalf("erreur = %v, attendue %v", err, testCase.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("erreur inattendue : %v", err)
			}
			if user.Pseudo != testCase.wantPseudo {
				t.Errorf("Pseudo = %q, attendu %q", user.Pseudo, testCase.wantPseudo)
			}
			if user.CreditBalance != 0 {
				t.Errorf("CreditBalance = %d, attendu 0 (les credits de bienvenue sont attribues par le use case)", user.CreditBalance)
			}
			if user.Skills == nil {
				t.Error("Skills doit etre une slice vide non nil")
			}
			if user.CreatedAt == "" {
				t.Error("CreatedAt ne doit pas etre vide")
			}
		})
	}
}
