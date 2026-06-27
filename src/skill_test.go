package main

import (
	"errors"
	"testing"
)

func TestNewSkill(t *testing.T) {
	t.Run("competence valide", func(t *testing.T) {
		skill, err := NewSkill("  Jardinage  ", "expert")
		if err != nil {
			t.Fatalf("erreur inattendue : %v", err)
		}
		if skill.Name != "Jardinage" {
			t.Errorf("Name = %q, attendu Jardinage (trim applique)", skill.Name)
		}
		if skill.Level != "expert" {
			t.Errorf("Level = %q, attendu expert", skill.Level)
		}
	})

	t.Run("name vide -> ErrSkillNameRequired", func(t *testing.T) {
		_, err := NewSkill("   ", "expert")
		if !errors.Is(err, ErrSkillNameRequired) {
			t.Fatalf("erreur = %v, attendue ErrSkillNameRequired", err)
		}
	})

	t.Run("level invalide -> ErrSkillLevelInvalid", func(t *testing.T) {
		_, err := NewSkill("Jardinage", "maitre")
		if !errors.Is(err, ErrSkillLevelInvalid) {
			t.Fatalf("erreur = %v, attendue ErrSkillLevelInvalid", err)
		}
	})

	t.Run("name hors enum -> ErrSkillNameInvalid", func(t *testing.T) {
		_, err := NewSkill("Cuisinier", "expert")
		if !errors.Is(err, ErrSkillNameInvalid) {
			t.Fatalf("erreur = %v, attendue ErrSkillNameInvalid", err)
		}
	})
}
