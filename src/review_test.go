package main

import (
	"errors"
	"testing"
)

func TestNewReview(t *testing.T) {
	testCases := []struct {
		name       string
		exchangeID int
		authorID   int
		targetID   int
		rating     int
		comment    string
		wantErr    error
	}{
		{
			name:       "valid review",
			exchangeID: 1,
			authorID:   2,
			targetID:   1,
			rating:     5,
			comment:    "  Excellent  ",
		},
		{
			name:       "rating too low",
			exchangeID: 1,
			authorID:   2,
			targetID:   1,
			rating:     0,
			wantErr:    ErrReviewRatingInvalid,
		},
		{
			name:       "rating too high",
			exchangeID: 1,
			authorID:   2,
			targetID:   1,
			rating:     6,
			wantErr:    ErrReviewRatingInvalid,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			review, err := NewReview(testCase.exchangeID, testCase.authorID, testCase.targetID, testCase.rating, testCase.comment)
			if testCase.wantErr != nil {
				if !errors.Is(err, testCase.wantErr) {
					t.Fatalf("error = %v, want %v", err, testCase.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if review.Comment != "Excellent" {
				t.Errorf("Comment = %q, want Excellent (trim applied)", review.Comment)
			}
			if review.CreatedAt == "" {
				t.Error("CreatedAt must not be empty")
			}
		})
	}
}
