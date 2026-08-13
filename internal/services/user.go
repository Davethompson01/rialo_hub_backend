package services

import (
	"github.com/Davethompson01/rialo_hub_backend/config"
	repository "github.com/Davethompson01/rialo_hub_backend/internal/Repository"
	"github.com/Davethompson01/rialo_hub_backend/internal/models"
)

func GetUserProfile(api *config.ApiConfig, userID int) (models.Profile, error) {
	profile, err := repository.SelectUserByID(api, userID)
	if err != nil {
		return models.Profile{}, err
	}

	return profile, nil
}
