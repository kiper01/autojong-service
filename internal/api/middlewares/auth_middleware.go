package middlewares

import (
	"net/http"
	"strings"
)

type Auth interface {
	IsValidToken(string) bool
}

func AuthMiddleware(next http.HandlerFunc, auth Auth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization scheme", http.StatusUnauthorized)
			return
		}
		token := parts[1]

		if !auth.IsValidToken(token) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}
