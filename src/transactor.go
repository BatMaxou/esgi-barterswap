package main

import (
	"context"
	"database/sql"
	"fmt"
)

type dbExecutor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

type Transactor struct {
	db *sql.DB
}

func NewTransactor(db *sql.DB) *Transactor {
	return &Transactor{db: db}
}

func (transactor *Transactor) Executor() dbExecutor {
	return transactor.db
}

func (transactor *Transactor) WithinTransaction(ctx context.Context, fn func(exec dbExecutor) error) error {
	tx, err := transactor.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("debut transaction : %w", err)
	}
	defer tx.Rollback()

	if err := fn(tx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit : %w", err)
	}

	return nil
}
