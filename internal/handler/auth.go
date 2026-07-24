package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Davethompson01/rialo_hub_backend/config"
	"github.com/Davethompson01/rialo_hub_backend/internal/models"
	"github.com/Davethompson01/rialo_hub_backend/internal/services"
)

func StudenthandlerCreateAccount(apicfg *config.ApiConfig) http.HandlerFunc {

	return func(res http.ResponseWriter, r *http.Request) {
		var register models.Register
		json.NewDecoder(r.Body).Decode(&register)

		studentServices, err := services.Register(apicfg, register)
		if err != nil {
			RespondWithJson(res, 400, false, err.Error(), nil)
			return
		}

		RespondWithJson(
			res,
			http.StatusCreated,
			true,
			studentServices,
			nil,
		)

	}
}
