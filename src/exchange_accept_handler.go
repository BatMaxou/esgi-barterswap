package main

import (
	"net/http"
	"strconv"
)

func (a *api) handleAcceptExchange(w http.ResponseWriter, r *http.Request) {
	actor, ok := currentUser(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "authentication required")

		return
	}

	exchangeID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")

		return
	}

	exchange, err := a.exchanges.Accept(r.Context(), actor.ID, exchangeID)
	if err != nil {
		writeExchangeError(w, err)

		return
	}

	writeJSON(w, http.StatusOK, exchange)
}
