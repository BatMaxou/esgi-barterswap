package main

// UserStats aggregates read-only metrics for a user profile dashboard.
// Values are computed by UserStatsUseCase from services, exchanges, reviews
// and the credit transaction ledger.
type UserStats struct {
	UserID             int     `json:"user_id"`
	ActiveServices     int     `json:"active_services"`
	CompletedExchanges int     `json:"completed_exchanges"`
	CreditBalance      int     `json:"credit_balance"`
	AverageRating      float64 `json:"average_rating"`
	ReviewCount        int     `json:"review_count"`
	TotalEarned        int     `json:"total_earned"`
	TotalSpent         int     `json:"total_spent"`
}
