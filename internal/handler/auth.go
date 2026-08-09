package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Davethompson01/rialo_hub_backend/config"
	auth "github.com/Davethompson01/rialo_hub_backend/internal/Auth"
	"github.com/Davethompson01/rialo_hub_backend/internal/models"
	"github.com/Davethompson01/rialo_hub_backend/internal/services"
)

func CreateUser(apicfg *config.ApiConfig) http.HandlerFunc {

	return func(res http.ResponseWriter, r *http.Request) {
		var register models.Register
		json.NewDecoder(r.Body).Decode(&register)

		create, err := services.Register(apicfg, register)
		if err != nil {
			RespondWithJson(res, 400, false, err.Error(), nil)
			return
		}

		RespondWithJson(
			res,
			http.StatusCreated,
			true,
			create,
			nil,
		)

	}
}

func Login(apicfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var loginModel models.Login
		json.NewDecoder(r.Body).Decode(&loginModel)

		login, err := services.LoginInto_AsStudent(apicfg, loginModel)
		if err != nil {
			RespondWithJson(w, 400, false, err.Error(), nil)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:     "access_token",
			Value:    login.AccessToken,
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
			MaxAge:   15 * 60,
		})

		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    login.RefreshToken,
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
			MaxAge:   7 * 24 * 60 * 60,
		})
		RespondWithJson(
			w,
			http.StatusCreated,
			true,
			"Login Successful",
			struct{}{},
		)
	}

}

func GetMe(api *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("refresh_token")
		if err != nil {
			RespondWithJson(w, http.StatusUnauthorized, false, "Not authenticated", nil)
			return
		}

		claims, err := auth.ValidateToken(cookie.Value)
		if err != nil {
			RespondWithJson(w, http.StatusUnauthorized, false, "Invalid or expired token", nil)
			return
		}

		RespondWithJson(
			w,
			http.StatusOK,
			true,
			"User fetched",
			map[string]interface{}{
				"id":   claims.UserID,
				"role": claims.Role,
			},
		)
	}
}
