package middleware

import (
	"context"
	"encoding/json"
	"net/http"

	auth "github.com/Davethompson01/rialo_hub_backend/internal/Auth"
)

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("access_token")
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)

			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"message": "Missing token",
				"data":    nil,
			})
			return
		}

		claims, err := auth.ValidateToken(cookie.Value)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)

			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"message": "Invalid or expired token",
				"data":    nil,
			})
			return
		}

		ctx := context.WithValue(
			r.Context(),
			ClaimsKey,
			claims,
		)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// func JWTMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

// 		authHeader := r.Header.Get("Authorization")

// 		if authHeader == "" {
// 			http.Error(w, "Missing token", http.StatusUnauthorized)
// 			return
// 		}

// 		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

// 		claims, err := auth.ValidateToken(tokenString)
// 		if err != nil {
// 			http.Error(w, fmt.Sprintf("Invalid token: %v", err), http.StatusUnauthorized)
// 			return
// 		}

// 		ctx := context.WithValue(r.Context(), ClaimsKey, claims)

// 		next.ServeHTTP(w, r.WithContext(ctx))
// 	})
// }
