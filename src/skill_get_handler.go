package main

import (
	"errors"
	"net/http"
	"strconv"
)

func (a *api) handleGetUserSkills(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid identifier")

		return
	}

	skills, err := a.skills.ListSkills(r.Context(), userID)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			writeError(w, http.StatusNotFound, err.Error())

			return
		}
		writeError(w, http.StatusInternalServerError, "could not fetch skills")

		return
	}

	writeJSON(w, http.StatusOK, skills)
}
