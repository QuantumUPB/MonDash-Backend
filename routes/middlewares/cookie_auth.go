package middlewares

import (
	"net/http"
	"os"
)

// CookieAuthMiddleware ensures requests contain a valid auth_token cookie.
func CookieAuthMiddleware(next http.Handler) http.Handler {
	requiredToken := os.Getenv("AUTH_TOKEN")
	if requiredToken == "" {
		requiredToken = "abc"
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("auth_token")
		if err != nil || cookie.Value == "" {
			http.Error(w, "invalid auth token", http.StatusUnauthorized)
			return
		}
		if cookie.Value != requiredToken {
			http.Error(w, "invalid auth token", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
