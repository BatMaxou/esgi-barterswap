package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

type ServiceRepository struct{}

func NewServiceRepository() *ServiceRepository {
	return &ServiceRepository{}
}

func (repository *ServiceRepository) Create(ctx context.Context, exec dbExecutor, service Service) (Service, error) {
	createdAt, err := time.Parse(time.RFC3339, service.CreatedAt)
	if err != nil {
		return Service{}, fmt.Errorf("invalid creation date: %w", err)
	}

	insertResult, err := exec.ExecContext(ctx,
		`INSERT INTO services (provider_id, title, description, category, duration_minutes, credits, city, active, created_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		service.ProviderID, service.Title, service.Description, service.Category,
		service.DurationMinutes, service.Credits, service.City, service.Active, createdAt,
	)
	if err != nil {
		return Service{}, fmt.Errorf("insert service: %w", err)
	}

	insertedID, err := insertResult.LastInsertId()
	if err != nil {
		return Service{}, fmt.Errorf("fetch inserted id: %w", err)
	}
	service.ID = int(insertedID)

	return service, nil
}

func (repository *ServiceRepository) Update(ctx context.Context, exec dbExecutor, service Service) (Service, error) {
	_, err := exec.ExecContext(ctx,
		`UPDATE services
		 SET title = ?, description = ?, category = ?, duration_minutes = ?, credits = ?, city = ?, active = ?
		 WHERE id = ?`,
		service.Title, service.Description, service.Category, service.DurationMinutes,
		service.Credits, service.City, service.Active, service.ID,
	)
	if err != nil {
		return Service{}, fmt.Errorf("update service: %w", err)
	}

	return service, nil
}

func (repository *ServiceRepository) Delete(ctx context.Context, exec dbExecutor, id int) error {
	if _, err := exec.ExecContext(ctx, `DELETE FROM services WHERE id = ?`, id); err != nil {
		return fmt.Errorf("delete service: %w", err)
	}

	return nil
}

func (repository *ServiceRepository) FindByID(ctx context.Context, exec dbExecutor, id int) (Service, error) {
	row := exec.QueryRowContext(ctx,
		`SELECT id, provider_id, title, description, category, duration_minutes, credits, city, active, created_at
		 FROM services WHERE id = ?`,
		id,
	)

	service, err := scanService(row)
	if errors.Is(err, sql.ErrNoRows) {
		return Service{}, ErrServiceNotFound
	}
	if err != nil {
		return Service{}, fmt.Errorf("fetch service: %w", err)
	}

	return service, nil
}

func (repository *ServiceRepository) List(ctx context.Context, exec dbExecutor, filter ServiceFilter) ([]Service, error) {
	query := `SELECT id, provider_id, title, description, category, duration_minutes, credits, city, active, created_at
		 FROM services`

	conditions := []string{}
	args := []any{}

	if filter.Category != "" {
		conditions = append(conditions, "category = ?")
		args = append(args, filter.Category)
	}
	if filter.City != "" {
		conditions = append(conditions, "city = ?")
		args = append(args, filter.City)
	}
	if filter.Search != "" {
		conditions = append(conditions, "(title LIKE ? OR description LIKE ?)")
		pattern := "%" + filter.Search + "%"
		args = append(args, pattern, pattern)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	query += " ORDER BY created_at DESC, id DESC"

	rows, err := exec.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("fetch services: %w", err)
	}
	defer rows.Close()

	services := []Service{}
	for rows.Next() {
		service, err := scanService(rows)
		if err != nil {
			return nil, fmt.Errorf("read service: %w", err)
		}
		services = append(services, service)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate services: %w", err)
	}

	return services, nil
}

// scanner abstracts reading a single row, whether it comes from a *sql.Row
// (FindByID) or a *sql.Rows (List).
type scanner interface {
	Scan(dest ...any) error
}

func scanService(row scanner) (Service, error) {
	var service Service
	var description sql.NullString
	var city sql.NullString
	var createdAt time.Time

	if err := row.Scan(
		&service.ID, &service.ProviderID, &service.Title, &description, &service.Category,
		&service.DurationMinutes, &service.Credits, &city, &service.Active, &createdAt,
	); err != nil {
		return Service{}, err
	}

	service.Description = description.String
	service.City = city.String
	service.CreatedAt = createdAt.UTC().Format(time.RFC3339)

	return service, nil
}
