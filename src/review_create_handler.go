package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type createReviewRequest struct {
	Rating  int    `json:"rating"`
	Comment string `json:"comment"`
}

func (a *api) handleCreateReview(w http.ResponseWriter, r *http.Request) {
	actor, ok := currentUser(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "authentication required")

		return
	}

	exchangeID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")

		return
	}

	var requestBody createReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON request body")

		return
	}

	created, err := a.reviews.Create(r.Context(), actor.ID, exchangeID, requestBody.Rating, requestBody.Comment)
	if err != nil {
		writeReviewError(w, err)

		return
	}

	writeJSON(w, http.StatusCreated, created)
}

func writeReviewError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrReviewRatingInvalid),
		errors.Is(err, ErrReviewExchangeNotCompleted),
		errors.Is(err, ErrReviewAlreadyExists):
		writeError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, ErrReviewNotParticipant):
		writeError(w, http.StatusForbidden, err.Error())
	case errors.Is(err, ErrExchangeNotFound),
		errors.Is(err, ErrUserNotFound),
		errors.Is(err, ErrServiceNotFound):
		writeError(w, http.StatusNotFound, err.Error())
	default:
		writeError(w, http.StatusInternalServerError, "could not process review")
	}
}
