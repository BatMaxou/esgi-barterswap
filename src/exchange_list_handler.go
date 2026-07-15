package main

import "net/http"

func (a *api) handleListExchanges(w http.ResponseWriter, r *http.Request) {
	actor, ok := currentUser(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "authentication required")

		return
	}

	exchanges, err := a.exchanges.List(r.Context(), actor.ID, r.URL.Query().Get("status"))
	if err != nil {
		writeExchangeError(w, err)

		return
	}
	if exchanges == nil {
		exchanges = []Exchange{}
	}

	writeJSON(w, http.StatusOK, exchanges)
}
