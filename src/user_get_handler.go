package main

import (
	"errors"
	"net/http"
	"strconv"
)

func (a *api) handleGetUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid identifier")

		return
	}

	user, err := a.users.GetProfile(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			writeError(w, http.StatusNotFound, err.Error())

			return
		}
		writeError(w, http.StatusInternalServerError, "could not fetch user")

		return
	}

	writeJSON(w, http.StatusOK, user)
}
