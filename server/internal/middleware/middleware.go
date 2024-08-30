package middleware

import (
	"context"
	"encoding/base64"
	"net/http"
	"strings"
)

type contextKey string

const (
	ContextKeyUserID     contextKey = "user-id"
	ContextKeyUserSecret contextKey = "user-secret"
)

func BasicAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessKey := r.Header.Get("Authorization")

		if accessKey == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if accessKey[:6] != "Basic " {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		accessKey = accessKey[6:]

		decoded, err := base64.StdEncoding.DecodeString(accessKey)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(string(decoded), ":")
		if len(parts) != 2 {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userID := parts[0]
		userSecret := parts[1]

		r = r.WithContext(context.WithValue(r.Context(), ContextKeyUserID, userID))
		r = r.WithContext(context.WithValue(r.Context(), ContextKeyUserSecret, userSecret))

		next.ServeHTTP(w, r)
	})
}
