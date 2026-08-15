package route

import (
	middleware "github.com/Davethompson01/rialo_hub_backend/Middleware"
	"github.com/Davethompson01/rialo_hub_backend/config"
	"github.com/Davethompson01/rialo_hub_backend/internal/handler"
	"github.com/go-chi/chi"
)

func Post(r chi.Router, api *config.ApiConfig) {
	r.Route("/posts", func(r chi.Router) {
		r.Use(middleware.APIKey)
		r.Use(middleware.JWTMiddleware)
		r.Use(middleware.RequireRole("admin", "super_admin", "Artist", "Writer", "Moderators"))

		r.Post("/create", handler.CreatePost(api))
		r.Get("/", handler.GetPostFeed(api))

		r.Get("/{postID}", handler.GetPost(api))
		r.Put("/{postID}", handler.UpdatePost(api))
		r.Delete("/{postID}", handler.DeletePost(api))

		r.Post("/{postID}/like", handler.LikePost(api))
		r.Delete("/{postID}/like", handler.UnlikePost(api))

		r.Post("/{postID}/comments", handler.CreateComment(api))
	})

	r.Route("/comments", func(r chi.Router) {
		r.Use(middleware.APIKey)
		r.Use(middleware.JWTMiddleware)
		r.Use(middleware.RequireRole("admin", "super_admin", "Artist", "Writer", "Moderators"))

		r.Put("/{commentID}", handler.UpdateComment(api))
		r.Delete("/{commentID}", handler.DeleteComment(api))
		r.Get("/{postID}/comments", handler.GetPostComments(api))

	})
}
