package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	// The database may not be ready immediately, retry for up to 30s.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for {
		pingErr := db.PingContext(ctx)
		if pingErr == nil {
			return db, nil
		}
		if ctx.Err() != nil {
			return nil, fmt.Errorf("connect to database: %w", pingErr)
		}
		time.Sleep(time.Second)
	}
}

func migrate(ctx context.Context, db *sql.DB) error {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id         INT AUTO_INCREMENT PRIMARY KEY,
			pseudo     VARCHAR(255) NOT NULL,
			bio        TEXT,
			city      VARCHAR(255),
			created_at DATETIME NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS skills (
			id      INT AUTO_INCREMENT PRIMARY KEY,
			user_id INT NOT NULL,
			name     VARCHAR(255) NOT NULL,
			level  VARCHAR(32) NOT NULL,
			CONSTRAINT fk_skill_user FOREIGN KEY (user_id) REFERENCES users(id)
		)`,
		`CREATE TABLE IF NOT EXISTS credit_transactions (
			id          INT AUTO_INCREMENT PRIMARY KEY,
			user_id     INT NOT NULL,
			exchange_id INT NULL,
			amount     INT NOT NULL,
			type        VARCHAR(32) NOT NULL,
			created_at  DATETIME NOT NULL,
			CONSTRAINT fk_credit_user FOREIGN KEY (user_id) REFERENCES users(id)
		)`,
		`CREATE TABLE IF NOT EXISTS services (
			id               INT AUTO_INCREMENT PRIMARY KEY,
			provider_id      INT NOT NULL,
			title            VARCHAR(255) NOT NULL,
			description      TEXT,
			category         VARCHAR(64) NOT NULL,
			duration_minutes INT NOT NULL,
			credits          INT NOT NULL,
			city             VARCHAR(255),
			active           BOOLEAN NOT NULL DEFAULT TRUE,
			created_at       DATETIME NOT NULL,
			CONSTRAINT fk_service_provider FOREIGN KEY (provider_id) REFERENCES users(id)
		)`,
		`CREATE TABLE IF NOT EXISTS exchanges (
			id           INT AUTO_INCREMENT PRIMARY KEY,
			service_id   INT NOT NULL,
			requester_id INT NOT NULL,
			owner_id     INT NOT NULL,
			status       VARCHAR(32) NOT NULL,
			created_at   DATETIME NOT NULL,
			updated_at   DATETIME NOT NULL,
			CONSTRAINT fk_exchange_service   FOREIGN KEY (service_id)   REFERENCES services(id),
			CONSTRAINT fk_exchange_requester FOREIGN KEY (requester_id) REFERENCES users(id),
			CONSTRAINT fk_exchange_owner     FOREIGN KEY (owner_id)     REFERENCES users(id)
		)`,
		`CREATE INDEX IF NOT EXISTS idx_exchanges_service_status ON exchanges (service_id, status)`,
	}

	for _, statement := range statements {
		if _, err := db.ExecContext(ctx, statement); err != nil {
			return fmt.Errorf("migration: %w", err)
		}
	}

	if err := ensureForeignKey(ctx, db, "credit_transactions", "fk_credit_exchange",
		`ALTER TABLE credit_transactions
		 ADD CONSTRAINT fk_credit_exchange FOREIGN KEY (exchange_id) REFERENCES exchanges(id)`,
	); err != nil {
		return fmt.Errorf("migration: %w", err)
	}

	return nil
}

func ensureForeignKey(ctx context.Context, db *sql.DB, table, name, statement string) error {
	var count int

	err := db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM information_schema.TABLE_CONSTRAINTS
		 WHERE CONSTRAINT_SCHEMA = DATABASE()
		   AND TABLE_NAME = ?
		   AND CONSTRAINT_NAME = ?
		   AND CONSTRAINT_TYPE = 'FOREIGN KEY'`,
		table, name,
	).Scan(&count)
	if err != nil {
		return fmt.Errorf("check foreign key %s: %w", name, err)
	}
	if count > 0 {
		return nil
	}

	if _, err := db.ExecContext(ctx, statement); err != nil {
		return fmt.Errorf("add foreign key %s: %w", name, err)
	}

	return nil
}
