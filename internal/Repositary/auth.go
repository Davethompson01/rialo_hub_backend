package repositary

import (
	"context"
	"time"

	"github.com/Davethompson01/rialo_hub_backend/config"
	"github.com/Davethompson01/rialo_hub_backend/internal/models"
)

func CreateUser(api *config.ApiConfig, models models.Register) error {

	query := `INSERT INTO users(discord_username, username, role, password, created_at)
		VALUES($1, $2, $3, $4, $5)`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := api.DB.ExecContext(
		ctx,
		query,
		models.DiscordUserName,
		models.UserName,
		models.Role,
		models.Password,
		models.Created_at,
	)

	return err
}

func CheckDiscordExist(apiCfg *config.ApiConfig, email string) bool {
	var exists bool

	query := `
		SELECT EXISTS (
			SELECT 1 FROM register WHERE email = $1
		)
	`

	err := apiCfg.DB.QueryRow(query, email).Scan(&exists)
	if err != nil {
		return false
	}

	return exists
}

func GetUserByUsername(apiCfg *config.ApiConfig, discord_username string) (models.Login, error) {
	var user models.Login

	query := `
		SELECT user_id, username, password, role
		FROM register
		WHERE discord_username = $1
		LIMIT 1;
	`

	err := apiCfg.DB.QueryRow(query, discord_username).Scan(
		&user.User_id,
		&user.Username,
		&user.Password,
		&user.Role,
	)

	if err != nil {
		return models.Login{}, err
	}

	return user, nil
}
