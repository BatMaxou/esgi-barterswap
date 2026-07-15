package main

import (
	"encoding/json"
	"errors"
	"net/http"
)

type createExchangeRequest struct {
	ServiceID int `json:"service_id"`
}

func (a *api) handleCreateExchange(w http.ResponseWriter, r *http.Request) {
	actor, ok := currentUser(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "authentication required")

		return
	}

	var requestBody createExchangeRequest
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON request body")

		return
	}

	created, err := a.exchanges.Create(r.Context(), actor.ID, requestBody.ServiceID)
	if err != nil {
		writeExchangeError(w, err)

		return
	}

	writeJSON(w, http.StatusCreated, created)
}

func writeExchangeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrExchangeServiceIDInvalid),
		errors.Is(err, ErrExchangeSelfRequest),
		errors.Is(err, ErrExchangeInsufficientCredits),
		errors.Is(err, ErrExchangeServiceInactive),
		errors.Is(err, ErrExchangeInvalidTransition),
		errors.Is(err, ErrExchangeStatusInvalid):
		writeError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, ErrForbidden):
		writeError(w, http.StatusForbidden, err.Error())
	case errors.Is(err, ErrExchangeNotFound),
		errors.Is(err, ErrServiceNotFound):
		writeError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, ErrExchangeServiceUnavailable):
		writeError(w, http.StatusConflict, err.Error())
	default:
		writeError(w, http.StatusInternalServerError, "could not process exchange")
	}
}
