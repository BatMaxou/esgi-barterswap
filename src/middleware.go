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
			writeError(w, http.StatusUnauthorized, "authentification requise (header X-User-ID)")

			return
		}

		id, err := strconv.Atoi(header)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "header X-User-ID invalide")

			return
		}

		user, err := a.users.Authenticate(r.Context(), id)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "utilisateur non authentifie")

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
