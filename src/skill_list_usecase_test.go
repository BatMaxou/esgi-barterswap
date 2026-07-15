package main

import (
	"context"
	"errors"
	"testing"
)

func TestSkillUseCaseListSkills(t *testing.T) {
	t.Run("existing user returns their skills", func(t *testing.T) {
		users := &fakeUserRepository{user: User{ID: 5}}
		skills := &fakeSkillRepository{skills: []Skill{{Name: "Jardinage", Level: "expert"}}}
		useCase := NewSkillUseCase(&fakeDatabase{}, users, skills)

		got, err := useCase.ListSkills(context.Background(), 5)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got) != 1 || got[0].Name != "Jardinage" {
			t.Errorf("skills = %+v, want [Jardinage]", got)
		}
	})

	t.Run("user not found -> ErrUserNotFound", func(t *testing.T) {
		users := &fakeUserRepository{findErr: ErrUserNotFound}
		skills := &fakeSkillRepository{}
		useCase := NewSkillUseCase(&fakeDatabase{}, users, skills)

		_, err := useCase.ListSkills(context.Background(), 999)
		if !errors.Is(err, ErrUserNotFound) {
			t.Fatalf("error = %v, want ErrUserNotFound", err)
		}
	})
}
