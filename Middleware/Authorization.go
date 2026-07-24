package middleware

import (
	"fmt"
	"net/http"

	auth "github.com/Davethompson01/rialo_hub_backend/internal/Auth"
)

type contextKey string

const ClaimsKey contextKey = "claims"

func RequireRole(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			value := r.Context().Value(ClaimsKey)
			if value == nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			claims, ok := value.(*auth.Claims)
			if !ok {
				http.Error(w, "Invalid token claims", http.StatusUnauthorized)
				return
			}

			userRole := claims.Role

			for _, allowed := range allowedRoles {
				if userRole == allowed {
					next.ServeHTTP(w, r)
					return
				}
			}
			fmt.Println("userID from token:", claims.UserID)
			fmt.Println("Allowed roles:", allowedRoles)

			http.Error(w, "Forbidden", http.StatusForbidden)
		})
	}
}
