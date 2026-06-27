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
		writeError(w, http.StatusBadRequest, "identifiant invalide")

		return
	}

	actor, ok := currentUser(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "authentification requise")

		return
	}

	var requestBody []Skill
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		writeError(w, http.StatusBadRequest, "corps de requete JSON invalide")

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
			writeError(w, http.StatusInternalServerError, "impossible de definir les competences")
		}

		return
	}

	writeJSON(w, http.StatusOK, skills)
}
