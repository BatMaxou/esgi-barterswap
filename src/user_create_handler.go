package main

import (
	"encoding/json"
	"errors"
	"net/http"
)

type createUserRequest struct {
	Pseudo string `json:"pseudo"`
	Bio    string `json:"bio"`
	City   string `json:"city"`
}

func (a *api) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	var requestBody createUserRequest
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON request body")

		return
	}

	created, err := a.users.Register(r.Context(), requestBody.Pseudo, requestBody.Bio, requestBody.City)
	if err != nil {
		if errors.Is(err, ErrPseudoRequired) {
			writeError(w, http.StatusBadRequest, err.Error())

			return
		}
		writeError(w, http.StatusInternalServerError, "could not create user")

		return
	}

	writeJSON(w, http.StatusCreated, created)
}
