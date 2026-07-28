package route

import (
	middleware "github.com/Davethompson01/rialo_hub_backend/Middleware"
	"github.com/Davethompson01/rialo_hub_backend/config"
	"github.com/Davethompson01/rialo_hub_backend/internal/handler"
	"github.com/go-chi/chi"
)

func CreateUser(r chi.Router, api *config.ApiConfig) {
	r.Route("/auth", func(r chi.Router) {
		r.Use(middleware.APIKey)

		// create user
		r.Post("/createuser", handler.CreateUser(api))

		//user login
		r.Post("/login", handler.Login(api))
	})
}
