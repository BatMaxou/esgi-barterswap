package main

import "context"

type fakeReviewRepository struct {
	createCalled    bool
	created         Review
	createErr       error
	userReviews     []Review
	exchangeReviews []Review
	listUserErr     error
	listExchangeErr error
}

func (fake *fakeReviewRepository) Create(ctx context.Context, exec dbExecutor, review Review) (Review, error) {
	fake.createCalled = true
	if fake.createErr != nil {
		return Review{}, fake.createErr
	}
	review.ID = 11
	fake.created = review
	return review, nil
}

func (fake *fakeReviewRepository) ListByTargetUserID(ctx context.Context, exec dbExecutor, userID int) ([]Review, error) {
	if fake.listUserErr != nil {
		return nil, fake.listUserErr
	}
	return fake.userReviews, nil
}

func (fake *fakeReviewRepository) ListByExchangeIDs(ctx context.Context, exec dbExecutor, exchangeIDs []int) ([]Review, error) {
	if fake.listExchangeErr != nil {
		return nil, fake.listExchangeErr
	}
	return fake.exchangeReviews, nil
}
