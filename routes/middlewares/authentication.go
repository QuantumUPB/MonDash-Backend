package middlewares

import (
	"net/http"
	"os"
	"strings"
)

// AuthMiddleware ensures requests contain X-Auth-Token header.
func AuthMiddleware(next http.Handler) http.Handler {
	requiredToken := os.Getenv("AUTH_TOKEN")
	if requiredToken == "" {
		requiredToken = "abc"
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("X-Auth-Token")
		if token == "" {
			http.Error(w, "invalid auth token", http.StatusUnauthorized)
			return
		}

		if strings.HasPrefix(strings.ToLower(token), "bearer ") {
			token = strings.TrimSpace(token[7:])
		}

		if token != requiredToken {
			http.Error(w, "invalid auth token", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
