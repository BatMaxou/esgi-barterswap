package main

import (
	"context"
	"errors"
	"testing"
)

type fakeSkillRepository struct {
	skills         []Skill
	findErr        error
	replaceCalled  bool
	replacedUserID int
	replacedSkills []Skill
	replaceErr     error
}

func (fake *fakeSkillRepository) FindByUserID(ctx context.Context, exec dbExecutor, userID int) ([]Skill, error) {
	if fake.findErr != nil {
		return nil, fake.findErr
	}
	return fake.skills, nil
}

func (fake *fakeSkillRepository) ReplaceForUser(ctx context.Context, exec dbExecutor, userID int, skills []Skill) error {
	fake.replaceCalled = true
	if fake.replaceErr != nil {
		return fake.replaceErr
	}
	fake.replacedUserID = userID
	fake.replacedSkills = skills
	return nil
}

func TestSkillUseCaseListSkills(t *testing.T) {
	t.Run("utilisateur existant renvoie ses competences", func(t *testing.T) {
		users := &fakeUserRepository{user: User{ID: 5}}
		skills := &fakeSkillRepository{skills: []Skill{{Name: "Jardinage", Level: "expert"}}}
		useCase := NewSkillUseCase(&fakeDatabase{}, users, skills)

		got, err := useCase.ListSkills(context.Background(), 5)
		if err != nil {
			t.Fatalf("erreur inattendue : %v", err)
		}
		if len(got) != 1 || got[0].Name != "Jardinage" {
			t.Errorf("competences = %+v, attendu [Jardinage]", got)
		}
	})

	t.Run("utilisateur introuvable -> ErrUserNotFound", func(t *testing.T) {
		users := &fakeUserRepository{findErr: ErrUserNotFound}
		skills := &fakeSkillRepository{}
		useCase := NewSkillUseCase(&fakeDatabase{}, users, skills)

		_, err := useCase.ListSkills(context.Background(), 999)
		if !errors.Is(err, ErrUserNotFound) {
			t.Fatalf("erreur = %v, attendue ErrUserNotFound", err)
		}
	})
}

func TestSkillUseCaseDefineSkills(t *testing.T) {
	t.Run("definition valide ecrase les competences", func(t *testing.T) {
		users := &fakeUserRepository{user: User{ID: 5}}
		skills := &fakeSkillRepository{}
		useCase := NewSkillUseCase(&fakeDatabase{}, users, skills)

		input := []Skill{{Name: "  Jardinage  ", Level: "expert"}, {Name: "Cuisine", Level: "débutant"}}
		got, err := useCase.DefineSkills(context.Background(), 5, 5, input)
		if err != nil {
			t.Fatalf("erreur inattendue : %v", err)
		}
		if !skills.replaceCalled {
			t.Error("ReplaceForUser doit etre appele")
		}
		if skills.replacedUserID != 5 {
			t.Errorf("userID remplace = %d, attendu 5", skills.replacedUserID)
		}
		if len(got) != 2 || got[0].Name != "Jardinage" {
			t.Errorf("competences = %+v, attendu name trimme", got)
		}
	})

	t.Run("competences d'un autre utilisateur -> ErrForbidden", func(t *testing.T) {
		users := &fakeUserRepository{user: User{ID: 5}}
		skills := &fakeSkillRepository{}
		useCase := NewSkillUseCase(&fakeDatabase{}, users, skills)

		_, err := useCase.DefineSkills(context.Background(), 5, 9, []Skill{{Name: "Jardinage", Level: "expert"}})
		if !errors.Is(err, ErrForbidden) {
			t.Fatalf("erreur = %v, attendue ErrForbidden", err)
		}
		if skills.replaceCalled {
			t.Error("aucune ecriture ne doit avoir lieu en cas d'acces interdit")
		}
	})

	t.Run("level invalide -> ErrSkillLevelInvalid", func(t *testing.T) {
		users := &fakeUserRepository{user: User{ID: 5}}
		skills := &fakeSkillRepository{}
		useCase := NewSkillUseCase(&fakeDatabase{}, users, skills)

		_, err := useCase.DefineSkills(context.Background(), 5, 5, []Skill{{Name: "Jardinage", Level: "maitre"}})
		if !errors.Is(err, ErrSkillLevelInvalid) {
			t.Fatalf("erreur = %v, attendue ErrSkillLevelInvalid", err)
		}
		if skills.replaceCalled {
			t.Error("aucune ecriture ne doit avoir lieu quand la validation echoue")
		}
	})

	t.Run("name vide -> ErrSkillNameRequired", func(t *testing.T) {
		users := &fakeUserRepository{user: User{ID: 5}}
		skills := &fakeSkillRepository{}
		useCase := NewSkillUseCase(&fakeDatabase{}, users, skills)

		_, err := useCase.DefineSkills(context.Background(), 5, 5, []Skill{{Name: "", Level: "expert"}})
		if !errors.Is(err, ErrSkillNameRequired) {
			t.Fatalf("erreur = %v, attendue ErrSkillNameRequired", err)
		}
		if skills.replaceCalled {
			t.Error("aucune ecriture ne doit avoir lieu quand la validation echoue")
		}
	})

	t.Run("name hors enum -> ErrSkillNameInvalid", func(t *testing.T) {
		users := &fakeUserRepository{user: User{ID: 5}}
		skills := &fakeSkillRepository{}
		useCase := NewSkillUseCase(&fakeDatabase{}, users, skills)

		_, err := useCase.DefineSkills(context.Background(), 5, 5, []Skill{{Name: "Cuisinier", Level: "expert"}})
		if !errors.Is(err, ErrSkillNameInvalid) {
			t.Fatalf("erreur = %v, attendue ErrSkillNameInvalid", err)
		}
		if skills.replaceCalled {
			t.Error("aucune ecriture ne doit avoir lieu quand la validation echoue")
		}
	})
}
