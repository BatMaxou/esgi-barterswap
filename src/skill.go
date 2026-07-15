package main

import (
	"errors"
	"strings"
)

var ErrSkillNameRequired = errors.New("skill name is required")

var ErrSkillLevelInvalid = errors.New("invalid skill level (débutant, intermédiaire or expert)")

var ErrSkillNameInvalid = errors.New("invalid skill name (not in the category list)")

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
