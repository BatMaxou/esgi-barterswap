package main

import (
	"net/http"
	"strconv"
)

func (a *api) handleListUserReviews(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid identifier")

		return
	}

	reviews, err := a.reviews.ListForUser(r.Context(), userID)
	if err != nil {
		writeReviewError(w, err)

		return
	}
	if reviews == nil {
		reviews = []Review{}
	}

	writeJSON(w, http.StatusOK, reviews)
}
