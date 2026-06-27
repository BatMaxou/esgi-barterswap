package main

import (
	"errors"
	"strings"
)

var ErrSkillNameRequired = errors.New("le nom de la competence est obligatoire")

var ErrSkillLevelInvalid = errors.New("niveau de competence invalide (debutant, intermediaire ou expert)")

var ErrSkillNameInvalid = errors.New("nom de competence invalide (hors de la liste des categories)")

type Skill struct {
	Name  string `json:"name"`
	Level string `json:"level"`
}

var validLevels = map[string]bool{
	"débutant":      true,
	"intermédiaire": true,
	"expert":        true,
}

var validCategories = map[string]bool{
	"Informatique": true,
	"Jardinage":    true,
	"Bricolage":    true,
	"Cuisine":      true,
	"Musique":      true,
	"Langues":      true,
	"Sport":        true,
	"Tutorat":      true,
	"Déménagement": true,
	"Photographie": true,
	"Animalier":    true,
	"Couture":      true,
	"Autre":        true,
}

func NewSkill(name, level string) (Skill, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return Skill{}, ErrSkillNameRequired
	}
	if !validCategories[name] {
		return Skill{}, ErrSkillNameInvalid
	}

	level = strings.TrimSpace(level)
	if !validLevels[level] {
		return Skill{}, ErrSkillLevelInvalid
	}

	return Skill{Name: name, Level: level}, nil
}
