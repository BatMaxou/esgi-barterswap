package main

import (
	"context"
	"fmt"
	"time"
)

type CreditTransactionRepository struct{}

func NewCreditTransactionRepository() *CreditTransactionRepository {
	return &CreditTransactionRepository{}
}

func (repository *CreditTransactionRepository) Create(ctx context.Context, exec dbExecutor, transaction CreditTransaction) error {
	createdAt, err := time.Parse(time.RFC3339, transaction.CreatedAt)
	if err != nil {
		return fmt.Errorf("date de creation invalide : %w", err)
	}

	var exchangeID any
	if transaction.ExchangeID != 0 {
		exchangeID = transaction.ExchangeID
	}

	if _, err := exec.ExecContext(ctx,
		`INSERT INTO credit_transactions (user_id, exchange_id, amount, type, created_at)
		 VALUES (?, ?, ?, ?, ?)`,
		transaction.UserID, exchangeID, transaction.Amount, transaction.Type, createdAt,
	); err != nil {
		return fmt.Errorf("insertion transaction de credits : %w", err)
	}

	return nil
}

func (repository *CreditTransactionRepository) BalanceByUserID(ctx context.Context, exec dbExecutor, userID int) (int, error) {
	var balance int

	err := exec.QueryRowContext(ctx,
		`SELECT COALESCE(SUM(amount), 0) FROM credit_transactions WHERE user_id = ?`,
		userID,
	).Scan(&balance)
	if err != nil {
		return 0, fmt.Errorf("calcul du solde : %w", err)
	}

	return balance, nil
}
