package main

import (
	"errors"
	"net/http"
	"strconv"
)

func (a *api) handleGetUserStats(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid identifier")

		return
	}

	stats, err := a.stats.Get(r.Context(), userID)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			writeError(w, http.StatusNotFound, err.Error())

			return
		}
		writeError(w, http.StatusInternalServerError, "could not fetch user stats")

		return
	}

	writeJSON(w, http.StatusOK, stats)
}
