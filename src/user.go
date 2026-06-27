package main

import (
	"errors"
	"strings"
	"time"
)

var ErrPseudoRequired = errors.New("le pseudo est obligatoire")

var ErrUserNotFound = errors.New("utilisateur introuvable")

var ErrForbidden = errors.New("action non autorisee")

type User struct {
	ID            int     `json:"id"`
	Pseudo        string  `json:"pseudo"`
	Bio           string  `json:"bio,omitempty"`
	City          string  `json:"city,omitempty"`
	Skills        []Skill `json:"skills,omitempty"`
	CreditBalance int     `json:"credit_balance"`
	CreatedAt     string  `json:"created_at"`
}

func NewUser(pseudo, bio, city string) (User, error) {
	pseudo = strings.TrimSpace(pseudo)
	if pseudo == "" {
		return User{}, ErrPseudoRequired
	}

	return User{
		Pseudo:    pseudo,
		Bio:       strings.TrimSpace(bio),
		City:      strings.TrimSpace(city),
		Skills:    []Skill{},
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}, nil
}
