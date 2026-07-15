package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

type ExchangeFilter struct {
	UserID int
	Status string
}

type ExchangeRepository struct{}

func NewExchangeRepository() *ExchangeRepository {
	return &ExchangeRepository{}
}

func (repository *ExchangeRepository) Create(ctx context.Context, exec dbExecutor, exchange Exchange) (Exchange, error) {
	createdAt, err := time.Parse(time.RFC3339, exchange.CreatedAt)
	if err != nil {
		return Exchange{}, fmt.Errorf("invalid creation date: %w", err)
	}
	updatedAt, err := time.Parse(time.RFC3339, exchange.UpdatedAt)
	if err != nil {
		return Exchange{}, fmt.Errorf("invalid update date: %w", err)
	}

	insertResult, err := exec.ExecContext(ctx,
		`INSERT INTO exchanges (service_id, requester_id, owner_id, status, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		exchange.ServiceID, exchange.RequesterID, exchange.OwnerID, exchange.Status, createdAt, updatedAt,
	)
	if err != nil {
		return Exchange{}, fmt.Errorf("insert exchange: %w", err)
	}

	insertedID, err := insertResult.LastInsertId()
	if err != nil {
		return Exchange{}, fmt.Errorf("fetch inserted id: %w", err)
	}
	exchange.ID = int(insertedID)

	return exchange, nil
}

func (repository *ExchangeRepository) FindByID(ctx context.Context, exec dbExecutor, id int) (Exchange, error) {
	row := exec.QueryRowContext(ctx,
		`SELECT id, service_id, requester_id, owner_id, status, created_at, updated_at
		 FROM exchanges WHERE id = ?`,
		id,
	)

	exchange, err := scanExchange(row)
	if errors.Is(err, sql.ErrNoRows) {
		return Exchange{}, ErrExchangeNotFound
	}
	if err != nil {
		return Exchange{}, fmt.Errorf("fetch exchange: %w", err)
	}

	return exchange, nil
}

func (repository *ExchangeRepository) UpdateStatus(ctx context.Context, exec dbExecutor, id int, status, updatedAt string) (Exchange, error) {
	parsedUpdatedAt, err := time.Parse(time.RFC3339, updatedAt)
	if err != nil {
		return Exchange{}, fmt.Errorf("invalid update date: %w", err)
	}

	updateResult, err := exec.ExecContext(ctx,
		`UPDATE exchanges SET status = ?, updated_at = ? WHERE id = ?`,
		status, parsedUpdatedAt, id,
	)
	if err != nil {
		return Exchange{}, fmt.Errorf("update exchange status: %w", err)
	}

	rowsAffected, err := updateResult.RowsAffected()
	if err != nil {
		return Exchange{}, fmt.Errorf("fetch rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return Exchange{}, ErrExchangeNotFound
	}

	return repository.FindByID(ctx, exec, id)
}

func (repository *ExchangeRepository) HasActiveForService(ctx context.Context, exec dbExecutor, serviceID int) (bool, error) {
	var exists bool

	err := exec.QueryRowContext(ctx,
		`SELECT EXISTS(
			SELECT 1 FROM exchanges
			WHERE service_id = ? AND status IN (?, ?)
		)`,
		serviceID, ExchangeStatusPending, ExchangeStatusAccepted,
	).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("check active exchange for service: %w", err)
	}

	return exists, nil
}

func (repository *ExchangeRepository) List(ctx context.Context, exec dbExecutor, filter ExchangeFilter) ([]Exchange, error) {
	query := `SELECT id, service_id, requester_id, owner_id, status, created_at, updated_at
		 FROM exchanges`

	conditions := []string{}
	args := []any{}

	if filter.UserID > 0 {
		conditions = append(conditions, "(requester_id = ? OR owner_id = ?)")
		args = append(args, filter.UserID, filter.UserID)
	}
	if filter.Status != "" {
		conditions = append(conditions, "status = ?")
		args = append(args, filter.Status)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	query += " ORDER BY updated_at DESC, id DESC"

	rows, err := exec.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("fetch exchanges: %w", err)
	}
	defer rows.Close()

	exchanges := []Exchange{}
	for rows.Next() {
		exchange, err := scanExchange(rows)
		if err != nil {
			return nil, fmt.Errorf("read exchange: %w", err)
		}
		exchanges = append(exchanges, exchange)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate exchanges: %w", err)
	}

	return exchanges, nil
}

func (repository *ExchangeRepository) ListByServiceID(ctx context.Context, exec dbExecutor, serviceID int) ([]Exchange, error) {
	rows, err := exec.QueryContext(ctx,
		`SELECT id, service_id, requester_id, owner_id, status, created_at, updated_at
		 FROM exchanges WHERE service_id = ? ORDER BY created_at DESC, id DESC`,
		serviceID,
	)
	if err != nil {
		return nil, fmt.Errorf("fetch exchanges for service: %w", err)
	}
	defer rows.Close()

	exchanges := []Exchange{}
	for rows.Next() {
		exchange, err := scanExchange(rows)
		if err != nil {
			return nil, fmt.Errorf("read exchange: %w", err)
		}
		exchanges = append(exchanges, exchange)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate exchanges for service: %w", err)
	}

	return exchanges, nil
}

func scanExchange(row scanner) (Exchange, error) {
	var exchange Exchange
	var createdAt time.Time
	var updatedAt time.Time

	if err := row.Scan(
		&exchange.ID, &exchange.ServiceID, &exchange.RequesterID, &exchange.OwnerID,
		&exchange.Status, &createdAt, &updatedAt,
	); err != nil {
		return Exchange{}, err
	}

	exchange.CreatedAt = createdAt.UTC().Format(time.RFC3339)
	exchange.UpdatedAt = updatedAt.UTC().Format(time.RFC3339)

	return exchange, nil
}
