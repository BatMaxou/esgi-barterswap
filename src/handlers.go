package main

import (
	"encoding/json"
	"net/http"
)

func registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /", handleIndex)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	user := NewUser("Thierry", "ma bio", "Paris")
	json.NewEncoder(w).Encode(user)
}
