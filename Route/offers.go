package route

import (
	middleware "github.com/Davethompson01/rialo_hub_backend/Middleware"
	"github.com/Davethompson01/rialo_hub_backend/config"
	"github.com/Davethompson01/rialo_hub_backend/internal/handler"
	"github.com/go-chi/chi"
)

func Offer(r chi.Router, api *config.ApiConfig) {
	r.Route("/offers", func(r chi.Router) {

		r.Use(middleware.APIKey)
		r.Use(middleware.JWTMiddleware)
		// r.Use(middleware.RequireRole("admin", "super_admin", "Artist", "Writer", "Moderators"))

		// Employer receives applicant offers
		r.Get(
			"/applicant",
			handler.GetApplicantOffersHandler(api),
		)
		// Applicant sees their own offers
		r.Get(
			"/my",
			handler.GetMyOffersHandler(api),
		)

		// Employer accepts offer
		r.Post(
			"/accept",
			handler.AcceptOfferHandler(api),
		)

		// Employer rejects offer
		r.Post(
			"/reject",
			handler.RejectOfferHandler(api),
		)
	}) 
}
