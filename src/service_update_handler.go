package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type updateServiceRequest struct {
	Title           string `json:"title"`
	Description     string `json:"description"`
	Category        string `json:"category"`
	DurationMinutes int    `json:"duration_minutes"`
	Credits         int    `json:"credits"`
	City            string `json:"city"`
	Active          *bool  `json:"active"`
}

func (a *api) handleUpdateService(w http.ResponseWriter, r *http.Request) {
	serviceID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid identifier")

		return
	}

	actor, ok := currentUser(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "authentication required")

		return
	}

	var requestBody updateServiceRequest
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON request body")

		return
	}

	updated, err := a.services.Update(r.Context(), actor.ID, serviceID,
		requestBody.Title, requestBody.Description, requestBody.Category, requestBody.City,
		requestBody.DurationMinutes, requestBody.Credits, requestBody.Active,
	)
	if err != nil {
		switch {
		case errors.Is(err, ErrForbidden):
			writeError(w, http.StatusForbidden, err.Error())
		case errors.Is(err, ErrServiceNotFound):
			writeError(w, http.StatusNotFound, err.Error())
		default:
			writeServiceValidationError(w, err)
		}

		return
	}

	writeJSON(w, http.StatusOK, updated)
}
