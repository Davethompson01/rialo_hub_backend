package services

import (
	"github.com/Davethompson01/rialo_hub_backend/config"
	repository "github.com/Davethompson01/rialo_hub_backend/internal/Repository"
	"github.com/Davethompson01/rialo_hub_backend/internal/models"
)

func GetDashboardFeed(
	api *config.ApiConfig,
	userID int,
) ([]models.DashboardFeed, error) {

	const limit = 15

	return repository.DashboardFeed(
		api,
		userID,
		limit,
	)
}
