package main

import "context"

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
