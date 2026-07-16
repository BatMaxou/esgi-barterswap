package main

import (
	"context"
	"log"
	"net/http"
	"time"
)

func main() {
	addr := ":" + getEnv("PORT", "8000")

	db, err := openDB(getEnv("DB_DSN", ""))
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := migrate(ctx, db); err != nil {
		log.Fatalf("migration: %v", err)
	}

	transactor := NewTransactor(db)
	userRepository := NewUserRepository()
	creditTransactionRepository := NewCreditTransactionRepository()
	skillRepository := NewSkillRepository()
	serviceRepository := NewServiceRepository()
	exchangeRepository := NewExchangeRepository()
	userUseCase := NewUserUseCase(transactor, userRepository, creditTransactionRepository)
	skillUseCase := NewSkillUseCase(transactor, userRepository, skillRepository)
	serviceUseCase := NewServiceUseCase(transactor, serviceRepository)
	exchangeUseCase := NewExchangeUseCase(transactor, exchangeRepository, serviceRepository, creditTransactionRepository)
	reviewRepository := NewReviewRepository()
	reviewUseCase := NewReviewUseCase(transactor, reviewRepository, exchangeRepository, userRepository, serviceRepository)
	userStatsUseCase := NewUserStatsUseCase(transactor, userRepository, serviceRepository, exchangeRepository, reviewRepository, creditTransactionRepository)
	app := &api{
		users:     userUseCase,
		skills:    skillUseCase,
		services:  serviceUseCase,
		exchanges: exchangeUseCase,
		reviews:   reviewUseCase,
		stats:     userStatsUseCase,
	}

	mux := http.NewServeMux()
	app.registerRoutes(mux)

	log.Printf("BarterSwap starting on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
