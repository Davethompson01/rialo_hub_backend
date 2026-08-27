package route

import (
	middleware "github.com/Davethompson01/rialo_hub_backend/Middleware"
	"github.com/Davethompson01/rialo_hub_backend/config"
	"github.com/Davethompson01/rialo_hub_backend/internal/handler"
	"github.com/go-chi/chi"
)

func Message(r chi.Router, api *config.ApiConfig) {
	r.Route("/conversations", func(r chi.Router) {
		r.Use(middleware.APIKey)
		r.Use(middleware.JWTMiddleware)
		r.Use(middleware.RequireRole("admin", "super_admin", "Artist", "Writer", "Moderators"))

		r.Post(
			"/message",
			handler.CreateMessage(api),
		)

		r.Post("/negotiate", handler.CreateNegotiation(api))

		r.Get(
			"/{conversationID}/messages",
			handler.GetConversationMessagesHandler(api),
		)

	})
}
