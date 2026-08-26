package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	middleware "github.com/Davethompson01/rialo_hub_backend/Middleware"
	"github.com/Davethompson01/rialo_hub_backend/config"
	auth "github.com/Davethompson01/rialo_hub_backend/internal/Auth"
	repository "github.com/Davethompson01/rialo_hub_backend/internal/Repository"
	"github.com/Davethompson01/rialo_hub_backend/internal/models"
	"github.com/Davethompson01/rialo_hub_backend/internal/services"
	"github.com/go-chi/chi"
)

func CreateNegotiation(api *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var message models.SendMessage
		if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
			RespondWithJson(
				w,
				http.StatusBadRequest,
				false,
				"Invalid request body",
				nil,
			)
			return
		}
		claims := r.Context().Value(middleware.ClaimsKey).(*auth.Claims)
		message.ApplicantID = claims.UserID
		message.CreatedBy = claims.UserID
		message.CreatedAt = time.Now()

		CreateNegotiation, err := services.CreateNegotiation(api, message)
		if err != nil {
			RespondWithJson(w, 400, false, err.Error(), nil)
			return
		}

		RespondWithJson(w, 200, true, "User can now Negotiate", CreateNegotiation)
	}
}

func GetApplicantOffersHandler(
	api *config.ApiConfig,
) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		claims, ok := r.Context().Value(middleware.ClaimsKey).(*auth.Claims)

		if !ok {
			http.Error(
				w,
				"unauthorized",
				http.StatusUnauthorized,
			)
			return
		}

		offers, err := services.GetAllApplicantOffer(
			api,
			claims.UserID,
		)

		if err != nil {
			RespondWithJson(w, http.StatusInternalServerError, false, err.Error(), nil)

			return
		}

		RespondWithJson(
			w,
			http.StatusOK,
			true,
			"Applicant fetched",
			offers,
		)
	}
}

func GetMyOffersHandler(
	api *config.ApiConfig,
) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		claims, ok := r.Context().Value(middleware.ClaimsKey).(*auth.Claims)

		if !ok {
			http.Error(
				w,
				"unauthorized",
				http.StatusUnauthorized,
			)
			return
		}

		offers, err := services.GetApplicationOffers(
			api,
			claims.UserID,
		)

		if err != nil {
			RespondWithJson(w, http.StatusInternalServerError, false, err.Error(), nil)

			return
		}

		RespondWithJson(
			w,
			http.StatusOK,
			true,
			"Offers fetched",
			offers,
		)
	}
}
func AcceptOfferHandler(
	api *config.ApiConfig,
) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		claims, ok := r.Context().Value(middleware.ClaimsKey).(*auth.Claims)

		if !ok {
			http.Error(
				w,
				"unauthorized",
				http.StatusUnauthorized,
			)
			return
		}

		var req models.OfferActionRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(
				w,
				"invalid request body",
				http.StatusBadRequest,
			)
			return
		}

		result, err := services.AcceptOffers(
			api,
			req.ApplicationID,
			req.TaskID,
			claims.UserID,
			req.OfferID,
			req.ConversationID,
		)

		if err != nil {
			RespondWithJson(
				w,
				http.StatusBadRequest,
				false,
				err.Error(),
				nil,
			)
			return
		}

		RespondWithJson(
			w,
			http.StatusCreated,
			true,
			"Accepted Offer",
			result,
		)
	}
}

func RejectOfferHandler(
	api *config.ApiConfig,
) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		claims, ok := r.Context().Value(middleware.ClaimsKey).(*auth.Claims)

		if !ok {
			http.Error(
				w,
				"unauthorized",
				http.StatusUnauthorized,
			)
			return
		}

		var req models.OfferActionRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(
				w,
				"invalid request body",
				http.StatusBadRequest,
			)
			return
		}

		result, err := services.RejectOffers(
			api,
			req.ApplicationID,
			req.TaskID,
			req.OfferID,
			claims.UserID,
			req.ConversationID,
		)

		if err != nil {
			RespondWithJson(
				w,
				http.StatusBadRequest,
				false,
				err.Error(),
				nil,
			)
			return
		}

		RespondWithJson(
			w,
			http.StatusCreated,
			true,
			"Rejected Offer",
			result,
		)
	}
}

func GetConversationMessagesHandler(
	api *config.ApiConfig,
) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		conversationID, err := strconv.Atoi(
			chi.URLParam(r, "conversationID"),
		)

		if err != nil {
			http.Error(
				w,
				"invalid conversation ID",
				http.StatusBadRequest,
			)
			return
		}

		messages, err := repository.GetMessages(
			api,
			conversationID,
		)

		if err != nil {
			RespondWithJson(
				w,
				http.StatusInternalServerError,
				false,
				err.Error(),
				nil,
			)
			return
		}

		RespondWithJson(
			w,
			http.StatusOK,
			true,
			"Messages fetched",
			messages,
		)
	}
}

