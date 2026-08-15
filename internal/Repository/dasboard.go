package repository

import (
	"context"
	"time"

	"github.com/Davethompson01/rialo_hub_backend/config"
	"github.com/Davethompson01/rialo_hub_backend/internal/models"
)

func DashboardFeed(
	api *config.ApiConfig,
	userID int,
	limit int,
) ([]models.DashboardFeed, error) {

	query := `
		SELECT
			id,
			type,
			user_id,
			username,
			title,
			description,
			created_at
		FROM (

			SELECT
				sp.post_id AS id,
				'post' AS type,
				u.user_id,
				u.username,
				sp.title,
				sp.description,
				sp.created_at
			FROM socialposts sp
			JOIN users u
				ON u.user_id = sp.user_id

			UNION ALL

			SELECT
				tk.task_id AS id,
				'task' AS type,
				u.user_id,
				u.username,
				tk.title,
				tk.description,
				tk.created_at
			FROM tasks tk
			JOIN users u
				ON u.user_id = tk.user_id

		) AS feed

		ORDER BY RANDOM()
		LIMIT $1
	`

	ctx, cancel := context.WithTimeout(
		context.Background(),
		3*time.Second,
	)
	defer cancel()

	rows, err := api.DB.QueryContext(
		ctx,
		query,
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feeds []models.DashboardFeed

	for rows.Next() {
		var feed models.DashboardFeed

		err := rows.Scan(
			&feed.ID,
			&feed.Type,
			&feed.UserID,
			&feed.Username,
			&feed.Title,
			&feed.Description,
			&feed.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		feeds = append(feeds, feed)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return feeds, nil
}
