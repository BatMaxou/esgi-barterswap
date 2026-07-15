package main

import "context"

// fakeSkillRepository is shared by the SkillUseCase operation tests.
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
