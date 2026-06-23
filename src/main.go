package main

import (
	"log"
	"net/http"
)

func main() {
	addr := ":" + getEnv("PORT", "8000")

	mux := http.NewServeMux()
	registerRoutes(mux)

	log.Printf("BarterSwap démarré sur %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("serveur arrêté : %v", err)
	}
}
