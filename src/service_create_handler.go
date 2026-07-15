package main

import (
	"encoding/json"
	"errors"
	"net/http"
)

type createServiceRequest struct {
	Title           string `json:"title"`
	Description     string `json:"description"`
	Category        string `json:"category"`
	DurationMinutes int    `json:"duration_minutes"`
	Credits         int    `json:"credits"`
	City            string `json:"city"`
}

func (a *api) handleCreateService(w http.ResponseWriter, r *http.Request) {
	actor, ok := currentUser(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "authentication required")

		return
	}

	var requestBody createServiceRequest
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON request body")

		return
	}

	created, err := a.services.Create(r.Context(), actor.ID,
		requestBody.Title, requestBody.Description, requestBody.Category, requestBody.City,
		requestBody.DurationMinutes, requestBody.Credits,
	)
	if err != nil {
		writeServiceValidationError(w, err)

		return
	}

	writeJSON(w, http.StatusCreated, created)
}

// writeServiceValidationError maps an ad's business errors to HTTP codes.
// Validation errors return 400, everything else 500.
func writeServiceValidationError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrServiceTitleRequired),
		errors.Is(err, ErrServiceCategoryInvalid),
		errors.Is(err, ErrServiceDurationInvalid),
		errors.Is(err, ErrServiceCreditsInvalid):
		writeError(w, http.StatusBadRequest, err.Error())
	default:
		writeError(w, http.StatusInternalServerError, "could not save service")
	}
}
