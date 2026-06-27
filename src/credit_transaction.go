package main

type CreditTransaction struct {
	ID         int    `json:"id"`
	UserID     int    `json:"user_id"`
	ExchangeID int    `json:"exchange_id"`
	Amount     int    `json:"amount"`
	Type       string `json:"type"`
	CreatedAt  string `json:"created_at"`
}
