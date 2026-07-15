package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type updateUserRequest struct {
	Pseudo string `json:"pseudo"`
	Bio    string `json:"bio"`
	City   string `json:"city"`
}

func (a *api) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
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

	var requestBody updateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON request body")

		return
	}

	updated, err := a.users.UpdateProfile(r.Context(), actor.ID, targetID, requestBody.Pseudo, requestBody.Bio, requestBody.City)
	if err != nil {
		switch {
		case errors.Is(err, ErrForbidden):
			writeError(w, http.StatusForbidden, err.Error())
		case errors.Is(err, ErrPseudoRequired):
			writeError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, ErrUserNotFound):
			writeError(w, http.StatusNotFound, err.Error())
		default:
			writeError(w, http.StatusInternalServerError, "could not update profile")
		}

		return
	}

	writeJSON(w, http.StatusOK, updated)
}
