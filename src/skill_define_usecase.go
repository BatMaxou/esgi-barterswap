package main

import (
	"context"
	"fmt"
)

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
		return nil, fmt.Errorf("set skills: %w", err)
	}

	return validated, nil
}
