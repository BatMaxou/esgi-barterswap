package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

func (a *api) handleDefineUserSkills(w http.ResponseWriter, r *http.Request) {
	targetID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid identifier")

		return
	}

	actor, ok := currentUser(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "authentication required")

		return
	}

	var requestBody []Skill
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON request body")

		return
	}

	skills, err := a.skills.DefineSkills(r.Context(), actor.ID, targetID, requestBody)
	if err != nil {
		switch {
		case errors.Is(err, ErrForbidden):
			writeError(w, http.StatusForbidden, err.Error())
		case errors.Is(err, ErrSkillNameRequired), errors.Is(err, ErrSkillNameInvalid), errors.Is(err, ErrSkillLevelInvalid):
			writeError(w, http.StatusBadRequest, err.Error())
		default:
			writeError(w, http.StatusInternalServerError, "could not set skills")
		}

		return
	}

	writeJSON(w, http.StatusOK, skills)
}
