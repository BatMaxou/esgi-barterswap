package main

import (
	"errors"
	"net/http"
	"strconv"
)

func (a *api) handleDeleteService(w http.ResponseWriter, r *http.Request) {
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

	if err := a.services.Delete(r.Context(), actor.ID, serviceID); err != nil {
		switch {
		case errors.Is(err, ErrForbidden):
			writeError(w, http.StatusForbidden, err.Error())
		case errors.Is(err, ErrServiceNotFound):
			writeError(w, http.StatusNotFound, err.Error())
		case errors.Is(err, ErrServiceHasExchanges):
			writeError(w, http.StatusConflict, err.Error())
		default:
			writeError(w, http.StatusInternalServerError, "could not delete service")
		}

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
