package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Davethompson01/rialo_hub_backend/config"
	"github.com/Davethompson01/rialo_hub_backend/internal/models"
	"github.com/Davethompson01/rialo_hub_backend/internal/services"
)

func CreateTasks(api *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var task models.Task
		json.NewDecoder(r.Body).Decode(&task)

		CreateTasks, err := services.CreateTasks(api, task)
		if err != nil {
			RespondWithJson(w, 400, false, err.Error(), nil)
			return
		}

		RespondWithJson(
			w,
			http.StatusCreated,
			true,
			CreateTasks,
			nil,
		)
	}

}
