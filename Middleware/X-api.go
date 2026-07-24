package middleware

import (
	"net/http"
	"os"
)

func APIKey(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		apiKey := r.Header.Get("X-API-Key")

		if apiKey == "" {
			http.Error(w, "Missing API Key", http.StatusUnauthorized)
			return
		}

		if apiKey != os.Getenv("API_KEY") {
			http.Error(w, "Invalid API Key", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}