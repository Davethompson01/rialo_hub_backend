package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Davethompson01/rialo_hub_backend/config"
	"github.com/Davethompson01/rialo_hub_backend/internal/models"
	"github.com/Davethompson01/rialo_hub_backend/internal/services"
)

func GetUserProfile(api *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var profile models.Profile
		
		if err := json.NewDecoder(r.Body).Decode(&profile); err != nil {
			RespondWithJson(w, http.StatusBadRequest, false, err.Error(), nil)
			return
		}

		fmt.Println(profile.UserID)
		getprofile, err := services.GetUserProfile(api, profile.UserID)
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
