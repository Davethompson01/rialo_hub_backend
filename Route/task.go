package route

import (
	middleware "github.com/Davethompson01/rialo_hub_backend/Middleware"
	"github.com/Davethompson01/rialo_hub_backend/config"
	"github.com/Davethompson01/rialo_hub_backend/internal/handler"
	"github.com/go-chi/chi"
)

func Task(r chi.Router, api *config.ApiConfig) {
	r.Route("/task", func(r chi.Router) {
		r.Use(middleware.APIKey)
		r.Use(middleware.JWTMiddleware)
		r.Use(middleware.RequireRole("admin", "super_admin", "Artist", "Writer"))

		r.Post("/create", handler.CreateTasks(api))
		r.Post("/apply", handler.ApplyForTasks(api))

	})
}
