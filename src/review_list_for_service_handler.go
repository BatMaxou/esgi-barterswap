package main

import (
	"net/http"
	"strconv"
)

func (a *api) handleListServiceReviews(w http.ResponseWriter, r *http.Request) {
	serviceID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid identifier")

		return
	}

	reviews, err := a.reviews.ListForService(r.Context(), serviceID)
	if err != nil {
		writeReviewError(w, err)

		return
	}
	if reviews == nil {
		reviews = []Review{}
	}

	writeJSON(w, http.StatusOK, reviews)
}
