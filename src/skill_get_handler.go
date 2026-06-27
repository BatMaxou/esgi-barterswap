package main

import (
	"errors"
	"net/http"
	"strconv"
)

func (a *api) handleGetUserSkills(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "identifiant invalide")

		return
	}

	skills, err := a.skills.ListSkills(r.Context(), userID)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			writeError(w, http.StatusNotFound, err.Error())

			return
		}
		writeError(w, http.StatusInternalServerError, "impossible de recuperer les competences")

		return
	}

	writeJSON(w, http.StatusOK, skills)
}
