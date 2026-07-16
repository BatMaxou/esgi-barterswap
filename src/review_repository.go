package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
)

type ReviewRepository struct{}

func NewReviewRepository() *ReviewRepository {
	return &ReviewRepository{}
}

func (repository *ReviewRepository) Create(ctx context.Context, exec dbExecutor, review Review) (Review, error) {
	createdAt, err := time.Parse(time.RFC3339, review.CreatedAt)
	if err != nil {
		return Review{}, fmt.Errorf("invalid creation date: %w", err)
	}

	insertResult, err := exec.ExecContext(ctx,
		`INSERT INTO reviews (exchange_id, author_id, target_id, rating, comment, created_at)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		review.ExchangeID, review.AuthorID, review.TargetID, review.Rating, review.Comment, createdAt,
	)
	if err != nil {
		var mysqlError *mysql.MySQLError
		if errors.As(err, &mysqlError) && mysqlError.Number == 1062 {
			return Review{}, ErrReviewAlreadyExists
		}
		return Review{}, fmt.Errorf("insert review: %w", err)
	}

	insertedID, err := insertResult.LastInsertId()
	if err != nil {
		return Review{}, fmt.Errorf("fetch inserted id: %w", err)
	}
	review.ID = int(insertedID)

	return review, nil
}

func (repository *ReviewRepository) ListByTargetUserID(ctx context.Context, exec dbExecutor, userID int) ([]Review, error) {
	rows, err := exec.QueryContext(ctx,
		`SELECT id, exchange_id, author_id, target_id, rating, comment, created_at
		 FROM reviews WHERE target_id = ? ORDER BY created_at DESC, id DESC`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("fetch reviews for user: %w", err)
	}
	defer rows.Close()

	return scanReviews(rows)
}

func (repository *ReviewRepository) StatsByTargetUserID(ctx context.Context, exec dbExecutor, userID int) (float64, int, error) {
	var averageRating float64
	var reviewCount int

	err := exec.QueryRowContext(ctx,
		`SELECT COALESCE(AVG(rating), 0), COUNT(*) FROM reviews WHERE target_id = ?`,
		userID,
	).Scan(&averageRating, &reviewCount)
	if err != nil {
		return 0, 0, fmt.Errorf("compute review stats: %w", err)
	}

	return averageRating, reviewCount, nil
}

func (repository *ReviewRepository) ListByExchangeIDs(ctx context.Context, exec dbExecutor, exchangeIDs []int) ([]Review, error) {
	if len(exchangeIDs) == 0 {
		return []Review{}, nil
	}

	placeholders := strings.Repeat("?,", len(exchangeIDs))
	placeholders = placeholders[:len(placeholders)-1]

	query := fmt.Sprintf(
		`SELECT id, exchange_id, author_id, target_id, rating, comment, created_at
		 FROM reviews WHERE exchange_id IN (%s) ORDER BY created_at DESC, id DESC`,
		placeholders,
	)

	args := make([]any, len(exchangeIDs))
	for index, exchangeID := range exchangeIDs {
		args[index] = exchangeID
	}

	rows, err := exec.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("fetch reviews for exchanges: %w", err)
	}
	defer rows.Close()

	return scanReviews(rows)
}

func scanReviews(rows *sql.Rows) ([]Review, error) {
	reviews := []Review{}
	for rows.Next() {
		review, err := scanReview(rows)
		if err != nil {
			return nil, fmt.Errorf("read review: %w", err)
		}
		reviews = append(reviews, review)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate reviews: %w", err)
	}

	return reviews, nil
}

func scanReview(row scanner) (Review, error) {
	var review Review
	var comment sql.NullString
	var createdAt time.Time

	if err := row.Scan(
		&review.ID, &review.ExchangeID, &review.AuthorID, &review.TargetID,
		&review.Rating, &comment, &createdAt,
	); err != nil {
		return Review{}, err
	}

	review.Comment = comment.String
	review.CreatedAt = createdAt.UTC().Format(time.RFC3339)

	return review, nil
}
