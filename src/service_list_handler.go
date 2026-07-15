package main

import "net/http"

func (a *api) handleListServices(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	filter := ServiceFilter{
		Category: query.Get("category"),
		City:     query.Get("city"),
		Search:   query.Get("search"),
	}

	services, err := a.services.List(r.Context(), filter)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not list services")

		return
	}

	writeJSON(w, http.StatusOK, services)
}
