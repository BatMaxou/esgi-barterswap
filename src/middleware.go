package main

import (
	"context"
	"net/http"
	"strconv"
)

type contextKey string

const currentUserKey contextKey = "currentUser"

func (a *api) requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("X-User-ID")
		if header == "" {
			writeError(w, http.StatusUnauthorized, "authentication required (X-User-ID header)")

			return
		}

		id, err := strconv.Atoi(header)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "invalid X-User-ID header")

			return
		}

		user, err := a.users.Authenticate(r.Context(), id)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "user not authenticated")

			return
		}

		ctx := context.WithValue(r.Context(), currentUserKey, user)

		next(w, r.WithContext(ctx))
	}
}

func currentUser(ctx context.Context) (User, bool) {
	user, ok := ctx.Value(currentUserKey).(User)

	return user, ok
}
