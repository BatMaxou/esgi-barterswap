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
		log.Fatalf("base de donnees : %v", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := migrate(ctx, db); err != nil {
		log.Fatalf("migration : %v", err)
	}

	transactor := NewTransactor(db)
	userRepository := NewUserRepository()
	creditTransactionRepository := NewCreditTransactionRepository()
	useCase := NewUserUseCase(transactor, userRepository, creditTransactionRepository)
	app := &api{
		users: useCase,
	}

	mux := http.NewServeMux()
	app.registerRoutes(mux)

	log.Printf("BarterSwap demarre sur %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("serveur arrete : %v", err)
	}
}
