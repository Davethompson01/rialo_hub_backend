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
		// r.Use(middleware.RequireRole("admin", "super_admin", "Artist", "Writer", "Moderators"))

		r.Post("/create", handler.CreateTasks(api))
		r.Post("/apply", handler.ApplyForTasks(api))
		r.Post("/reject", handler.RejectEmployee(api))
		r.Post("/accept", handler.AcceptEmployee(api))
		r.Get("/{taskID}/getTaskApplications", handler.GetTaskApplications(api))
		r.Get("/getMyApplications", handler.GetMyApplications(api))
		r.Post("/cancel", handler.CancelApplication(api))
		r.Post("/delete", handler.DeleteTask(api))
		r.Get("/taskfeeds", handler.Taskfeed(api))
	})
}
