package main

import (
	"errors"
	"strings"
	"time"
)

var ErrReviewNotFound = errors.New("review not found")

var ErrReviewRatingInvalid = errors.New("rating must be between 1 and 5")

var ErrReviewExchangeNotCompleted = errors.New("exchange must be completed before reviewing")

var ErrReviewAlreadyExists = errors.New("review already submitted for this exchange")

var ErrReviewNotParticipant = errors.New("only exchange participants can submit a review")

type Review struct {
	ID         int    `json:"id"`
	ExchangeID int    `json:"exchange_id"`
	AuthorID   int    `json:"author_id"`
	TargetID   int    `json:"target_id"`
	Rating     int    `json:"rating"`
	Comment    string `json:"comment,omitempty"`
	CreatedAt  string `json:"created_at"`
}

// NewReview validates the fields and returns a Review ready to be persisted.
// Exchange status and participant checks belong in the use case.
func NewReview(exchangeID, authorID, targetID, rating int, comment string) (Review, error) {
	if rating < 1 || rating > 5 {
		return Review{}, ErrReviewRatingInvalid
	}

	return Review{
		ExchangeID: exchangeID,
		AuthorID:   authorID,
		TargetID:   targetID,
		Rating:     rating,
		Comment:    strings.TrimSpace(comment),
		CreatedAt:  time.Now().UTC().Format(time.RFC3339),
	}, nil
}
