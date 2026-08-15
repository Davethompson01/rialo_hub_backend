package route

import (
	middleware "github.com/Davethompson01/rialo_hub_backend/Middleware"
	"github.com/Davethompson01/rialo_hub_backend/config"
	"github.com/Davethompson01/rialo_hub_backend/internal/handler"
	"github.com/go-chi/chi"
)

func Dashboard(r chi.Router, api *config.ApiConfig) {

	r.Route("/dashboard", func(r chi.Router) {
		r.Use(middleware.APIKey)
		r.Use(middleware.JWTMiddleware)
		r.Use(middleware.RequireRole("admin", "super_admin", "Artist", "Writer", "Moderators"))
		r.Get("/feed", handler.DashboardFeed(api))
	})
}
