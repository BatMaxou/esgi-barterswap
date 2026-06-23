package main

import "time"

const welcomeCredits = 10

type Skill struct {
	Nom    string `json:"nom"`
	Niveau string `json:"niveau"`
}

type User struct {
	ID            int     `json:"id"`
	Pseudo        string  `json:"pseudo"`
	Bio           string  `json:"bio,omitempty"`
	Ville         string  `json:"ville,omitempty"`
	Skills        []Skill `json:"skills,omitempty"`
	CreditBalance int     `json:"credit_balance"`
	CreatedAt     string  `json:"created_at"`
}

func NewUser(pseudo, bio, ville string) User {
	return User{
		Pseudo:        pseudo,
		Bio:           bio,
		Ville:         ville,
		Skills:        []Skill{},
		CreditBalance: welcomeCredits,
		CreatedAt:     time.Now().Format(time.RFC3339),
	}
}
