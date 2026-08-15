package services

import (
	"fmt"

	"github.com/Davethompson01/rialo_hub_backend/config"
	repository "github.com/Davethompson01/rialo_hub_backend/internal/Repository"
	"github.com/Davethompson01/rialo_hub_backend/internal/models"
)

func GetDashboardFeed(
	api *config.ApiConfig,
	userID int,
) ([]models.DashboardFeed, error) {

	feeds, err := repository.DashboardFeed(api, userID, 15)
	if err != nil {
		return nil, fmt.Errorf("failed to get dashboard feed: %w", err)
	}

	return feeds, nil
}
