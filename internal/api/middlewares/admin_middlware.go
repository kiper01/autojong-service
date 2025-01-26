package middlewares

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"
)

func AdminMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		parts := strings.Split(authHeader, " ")
		token := parts[1]

		tokenParts := strings.Split(token, ".")

		payloadBytes, err := base64.RawURLEncoding.DecodeString(tokenParts[1])
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var payload struct {
			Id    string   `json:"sub"`
			Roles []string `json:"roles"`
		}

		err = json.Unmarshal(payloadBytes, &payload)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if !hasAdminRole(payload.Roles) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		next(w, r)
	}
}

func hasAdminRole(roles []string) bool {
	for _, role := range roles {
		if role == "ROLE_ADMIN" {
			return true
		}
	}
	return false
}
