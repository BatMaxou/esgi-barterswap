package main

import (
	"context"
	"errors"
	"testing"
)

func TestSkillUseCaseDefineSkills(t *testing.T) {
	t.Run("valid definition overwrites the skills", func(t *testing.T) {
		users := &fakeUserRepository{user: User{ID: 5}}
		skills := &fakeSkillRepository{}
		useCase := NewSkillUseCase(&fakeDatabase{}, users, skills)

		input := []Skill{{Name: "  Jardinage  ", Level: "expert"}, {Name: "Cuisine", Level: "débutant"}}
		got, err := useCase.DefineSkills(context.Background(), 5, 5, input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !skills.replaceCalled {
			t.Error("ReplaceForUser must be called")
		}
		if skills.replacedUserID != 5 {
			t.Errorf("replaced userID = %d, want 5", skills.replacedUserID)
		}
		if len(got) != 2 || got[0].Name != "Jardinage" {
			t.Errorf("skills = %+v, want trimmed name", got)
		}
	})

	t.Run("another user's skills -> ErrForbidden", func(t *testing.T) {
		users := &fakeUserRepository{user: User{ID: 5}}
		skills := &fakeSkillRepository{}
		useCase := NewSkillUseCase(&fakeDatabase{}, users, skills)

		_, err := useCase.DefineSkills(context.Background(), 5, 9, []Skill{{Name: "Jardinage", Level: "expert"}})
		if !errors.Is(err, ErrForbidden) {
			t.Fatalf("error = %v, want ErrForbidden", err)
		}
		if skills.replaceCalled {
			t.Error("no write must happen on a forbidden access")
		}
	})

	t.Run("invalid level -> ErrSkillLevelInvalid", func(t *testing.T) {
		users := &fakeUserRepository{user: User{ID: 5}}
		skills := &fakeSkillRepository{}
		useCase := NewSkillUseCase(&fakeDatabase{}, users, skills)

		_, err := useCase.DefineSkills(context.Background(), 5, 5, []Skill{{Name: "Jardinage", Level: "maitre"}})
		if !errors.Is(err, ErrSkillLevelInvalid) {
			t.Fatalf("error = %v, want ErrSkillLevelInvalid", err)
		}
		if skills.replaceCalled {
			t.Error("no write must happen when validation fails")
		}
	})

	t.Run("empty name -> ErrSkillNameRequired", func(t *testing.T) {
		users := &fakeUserRepository{user: User{ID: 5}}
		skills := &fakeSkillRepository{}
		useCase := NewSkillUseCase(&fakeDatabase{}, users, skills)

		_, err := useCase.DefineSkills(context.Background(), 5, 5, []Skill{{Name: "", Level: "expert"}})
		if !errors.Is(err, ErrSkillNameRequired) {
			t.Fatalf("error = %v, want ErrSkillNameRequired", err)
		}
		if skills.replaceCalled {
			t.Error("no write must happen when validation fails")
		}
	})

	t.Run("name outside the enum -> ErrSkillNameInvalid", func(t *testing.T) {
		users := &fakeUserRepository{user: User{ID: 5}}
		skills := &fakeSkillRepository{}
		useCase := NewSkillUseCase(&fakeDatabase{}, users, skills)

		_, err := useCase.DefineSkills(context.Background(), 5, 5, []Skill{{Name: "Cuisinier", Level: "expert"}})
		if !errors.Is(err, ErrSkillNameInvalid) {
			t.Fatalf("error = %v, want ErrSkillNameInvalid", err)
		}
		if skills.replaceCalled {
			t.Error("no write must happen when validation fails")
		}
	})
}
