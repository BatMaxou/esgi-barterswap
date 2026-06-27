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
		return nil, fmt.Errorf("ouverture base : %w", err)
	}

	// La base peut ne pas etre prete immediatement, reessaie jusqu'a 30s
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for {
		pingErr := db.PingContext(ctx)
		if pingErr == nil {
			return db, nil
		}
		if ctx.Err() != nil {
			return nil, fmt.Errorf("connexion base : %w", pingErr)
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
			ville      VARCHAR(255),
			created_at DATETIME NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS credit_transactions (
			id          INT AUTO_INCREMENT PRIMARY KEY,
			user_id     INT NOT NULL,
			exchange_id INT NULL,
			montant     INT NOT NULL,
			type        VARCHAR(32) NOT NULL,
			created_at  DATETIME NOT NULL,
			CONSTRAINT fk_credit_user FOREIGN KEY (user_id) REFERENCES users(id)
		)`,
	}

	for _, statement := range statements {
		if _, err := db.ExecContext(ctx, statement); err != nil {
			return fmt.Errorf("migration : %w", err)
		}
	}

	return nil
}
