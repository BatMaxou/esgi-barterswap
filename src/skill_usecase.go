package main

import (
	"context"
	"fmt"
)

type skillRepository interface {
	FindByUserID(ctx context.Context, exec dbExecutor, userID int) ([]Skill, error)
	ReplaceForUser(ctx context.Context, exec dbExecutor, userID int, skills []Skill) error
}

type SkillUseCase struct {
	db     database
	users  userRepository
	skills skillRepository
}

func NewSkillUseCase(db database, users userRepository, skills skillRepository) *SkillUseCase {
	return &SkillUseCase{
		db:     db,
		users:  users,
		skills: skills,
	}
}

func (useCase *SkillUseCase) ListSkills(ctx context.Context, userID int) ([]Skill, error) {
	exec := useCase.db.Executor()

	if _, err := useCase.users.FindByID(ctx, exec, userID); err != nil {
		return nil, err
	}

	return useCase.skills.FindByUserID(ctx, exec, userID)
}

func (useCase *SkillUseCase) DefineSkills(ctx context.Context, actorID, targetID int, skills []Skill) ([]Skill, error) {
	if actorID != targetID {
		return nil, ErrForbidden
	}

	validated := make([]Skill, 0, len(skills))
	for _, skill := range skills {
		valid, err := NewSkill(skill.Name, skill.Level)
		if err != nil {
			return nil, err
		}
		validated = append(validated, valid)
	}

	err := useCase.db.WithinTransaction(ctx, func(exec dbExecutor) error {
		return useCase.skills.ReplaceForUser(ctx, exec, targetID, validated)
	})
	if err != nil {
		return nil, fmt.Errorf("definition des competences : %w", err)
	}

	return validated, nil
}
