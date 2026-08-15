package handler

import (
	"net/http"

	middleware "github.com/Davethompson01/rialo_hub_backend/Middleware"
	"github.com/Davethompson01/rialo_hub_backend/config"
	auth "github.com/Davethompson01/rialo_hub_backend/internal/Auth"
	"github.com/Davethompson01/rialo_hub_backend/internal/services"
)

func DashboardFeed(api *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		claims, ok := r.Context().Value(
			middleware.ClaimsKey,
		).(*auth.Claims)

		if !ok || claims == nil {
			RespondWithJson(
				w,
				http.StatusUnauthorized,
				false,
				"Unauthorized",
				nil,
			)
			return
		}

		// posts and tasks.
		feeds, err := services.GetDashboardFeed(
			api,
			claims.UserID,
		)

		if err != nil {
			RespondWithJson(
				w,
				http.StatusInternalServerError,
				false,
				err.Error(),
				nil,
			)
			return
		}

		RespondWithJson(
			w,
			http.StatusOK,
			true,
			"Dashboard feed retrieved successfully",
			feeds,
		)
	}
}
