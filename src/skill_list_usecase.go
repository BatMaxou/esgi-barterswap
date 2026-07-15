package main

import "context"

func (useCase *SkillUseCase) ListSkills(ctx context.Context, userID int) ([]Skill, error) {
	exec := useCase.db.Executor()

	if _, err := useCase.users.FindByID(ctx, exec, userID); err != nil {
		return nil, err
	}

	return useCase.skills.FindByUserID(ctx, exec, userID)
}
