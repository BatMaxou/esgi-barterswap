package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (repository *UserRepository) Create(ctx context.Context, exec dbExecutor, user User) (User, error) {
	createdAt, err := time.Parse(time.RFC3339, user.CreatedAt)
	if err != nil {
		return User{}, fmt.Errorf("date de creation invalide : %w", err)
	}

	insertResult, err := exec.ExecContext(ctx,
		`INSERT INTO users (pseudo, bio, city, created_at) VALUES (?, ?, ?, ?)`,
		user.Pseudo, user.Bio, user.City, createdAt,
	)
	if err != nil {
		return User{}, fmt.Errorf("insertion utilisateur : %w", err)
	}

	insertedID, err := insertResult.LastInsertId()
	if err != nil {
		return User{}, fmt.Errorf("recuperation de l'id : %w", err)
	}
	user.ID = int(insertedID)

	return user, nil
}

func (repository *UserRepository) Update(ctx context.Context, exec dbExecutor, user User) (User, error) {
	_, err := exec.ExecContext(ctx,
		`UPDATE users SET pseudo = ?, bio = ?, city = ? WHERE id = ?`,
		user.Pseudo, user.Bio, user.City, user.ID,
	)

	if err != nil {
		return User{}, fmt.Errorf("mise a jour utilisateur : %w", err)
	}

	return user, nil
}

func (repository *UserRepository) FindByID(ctx context.Context, exec dbExecutor, id int) (User, error) {
	var user User
	var createdAt time.Time

	err := exec.QueryRowContext(ctx,
		`SELECT id, pseudo, bio, city, created_at FROM users WHERE id = ?`,
		id,
	).Scan(&user.ID, &user.Pseudo, &user.Bio, &user.City, &createdAt)
	if errors.Is(err, sql.ErrNoRows) {
		return User{}, ErrUserNotFound
	}
	if err != nil {
		return User{}, fmt.Errorf("recuperation utilisateur : %w", err)
	}

	user.Skills = []Skill{}
	user.CreatedAt = createdAt.UTC().Format(time.RFC3339)
	return user, nil
}
