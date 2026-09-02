package handler

import (
	"net/http"

	middleware "github.com/Davethompson01/rialo_hub_backend/Middleware"
	"github.com/Davethompson01/rialo_hub_backend/config"
	auth "github.com/Davethompson01/rialo_hub_backend/internal/Auth"
	"github.com/Davethompson01/rialo_hub_backend/internal/services"
)

func GetUserProfile(api *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		claimsValue := r.Context().Value(middleware.ClaimsKey)

		claims, ok := claimsValue.(*auth.Claims)
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

		getprofile, err := services.GetUserProfile(api, claims.UserID)
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
			"Profile Fetched Successfully",
			getprofile,
		)
	}
}
