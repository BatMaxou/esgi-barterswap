package main

import (
	"errors"
	"net/http"
	"strconv"
)

func (a *api) handleGetService(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid identifier")

		return
	}

	service, err := a.services.Get(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrServiceNotFound) {
			writeError(w, http.StatusNotFound, err.Error())

			return
		}
		writeError(w, http.StatusInternalServerError, "could not fetch service")

		return
	}

	writeJSON(w, http.StatusOK, service)
}
