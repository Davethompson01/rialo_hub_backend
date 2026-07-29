package repository

import (
	"github.com/Davethompson01/rialo_hub_backend/config"
	"github.com/Davethompson01/rialo_hub_backend/internal/models"
)

func SelectUserByID(api *config.ApiConfig, user_id int) (models.Profile, error) {
	var profile models.Profile
	query := `SELECT
    user_id,
    COALESCE(profile_pics, '') AS profile_pics,
    discord_username,
    username,
    role
FROM users
WHERE user_id = $1`
	err := api.DB.QueryRow(query, user_id).Scan(
		&profile.UserID,
		&profile.Profile_pics,
		&profile.Discord_username,
		&profile.Username,
		&profile.Role,
	)
	if err != nil {
		return models.Profile{}, err
	}

	return profile, nil
}
