package main

import (
	"errors"
	"testing"
)

func TestNewSkill(t *testing.T) {
	t.Run("valid skill", func(t *testing.T) {
		skill, err := NewSkill("  Jardinage  ", "expert")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if skill.Name != "Jardinage" {
			t.Errorf("Name = %q, want Jardinage (trim applied)", skill.Name)
		}
		if skill.Level != "expert" {
			t.Errorf("Level = %q, want expert", skill.Level)
		}
	})

	t.Run("empty name -> ErrSkillNameRequired", func(t *testing.T) {
		_, err := NewSkill("   ", "expert")
		if !errors.Is(err, ErrSkillNameRequired) {
			t.Fatalf("error = %v, want ErrSkillNameRequired", err)
		}
	})

	t.Run("invalid level -> ErrSkillLevelInvalid", func(t *testing.T) {
		_, err := NewSkill("Jardinage", "maitre")
		if !errors.Is(err, ErrSkillLevelInvalid) {
			t.Fatalf("error = %v, want ErrSkillLevelInvalid", err)
		}
	})

	t.Run("name outside the enum -> ErrSkillNameInvalid", func(t *testing.T) {
		_, err := NewSkill("Cuisinier", "expert")
		if !errors.Is(err, ErrSkillNameInvalid) {
			t.Fatalf("error = %v, want ErrSkillNameInvalid", err)
		}
	})
}
