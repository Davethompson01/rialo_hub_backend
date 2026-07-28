package handler

import (
	"encoding/json"
	"net/http"

	middleware "github.com/Davethompson01/rialo_hub_backend/Middleware"
	"github.com/Davethompson01/rialo_hub_backend/config"
	auth "github.com/Davethompson01/rialo_hub_backend/internal/Auth"
	"github.com/Davethompson01/rialo_hub_backend/internal/models"
	"github.com/Davethompson01/rialo_hub_backend/internal/services"
)

func CreateTasks(api *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var task models.Task
		json.NewDecoder(r.Body).Decode(&task)

		claims := r.Context().Value(middleware.ClaimsKey).(*auth.Claims)
		task.UserID = claims.UserID

		CreateTasks, err := services.CreateTasks(api, task)
		if err != nil {
			RespondWithJson(w, 400, false, err.Error(), nil)
			return
		}

		RespondWithJson(
			w,
			200,
			true,
			CreateTasks,
			nil,
		)
	}

}
