package main

import "net/http"

func handleIndex(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"service": "BarterSwap",
		"status":  "ok",
	})
}
