package main

import (
	"context"
	"fmt"
)

type SkillRepository struct{}

func NewSkillRepository() *SkillRepository {
	return &SkillRepository{}
}

func (repository *SkillRepository) FindByUserID(ctx context.Context, exec dbExecutor, userID int) ([]Skill, error) {
	rows, err := exec.QueryContext(ctx,
		`SELECT name, level FROM skills WHERE user_id = ?`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("fetch skills: %w", err)
	}
	defer rows.Close()

	skills := []Skill{}
	for rows.Next() {
		var skill Skill
		if err := rows.Scan(&skill.Name, &skill.Level); err != nil {
			return nil, fmt.Errorf("read skill: %w", err)
		}
		skills = append(skills, skill)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate skills: %w", err)
	}

	return skills, nil
}

func (repository *SkillRepository) ReplaceForUser(ctx context.Context, exec dbExecutor, userID int, skills []Skill) error {
	if _, err := exec.ExecContext(ctx,
		`DELETE FROM skills WHERE user_id = ?`,
		userID,
	); err != nil {
		return fmt.Errorf("delete skills: %w", err)
	}

	for _, skill := range skills {
		if _, err := exec.ExecContext(ctx,
			`INSERT INTO skills (user_id, name, level) VALUES (?, ?, ?)`,
			userID, skill.Name, skill.Level,
		); err != nil {
			return fmt.Errorf("insert skill: %w", err)
		}
	}

	return nil
}
