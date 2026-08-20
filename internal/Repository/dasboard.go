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
		likes,
		comments,
		is_liked,
		created_at
	FROM (

		-- POSTS
		SELECT
			sp.post_id AS id,
			'post' AS type,
			u.user_id,
			u.username,
			sp.title,
			sp.description,
			sp.likes,
			sp.comments,
			EXISTS (
				SELECT 1
				FROM likes l
				WHERE l.post_id = sp.post_id
				  AND l.user_id = $1
			) AS is_liked,
			sp.created_at
		FROM socialposts sp
		JOIN users u
			ON u.user_id = sp.user_id

		UNION ALL

		-- TASKS
		SELECT
			tk.task_id AS id,
			'task' AS type,
			u.user_id,
			u.username,
			tk.title,
			tk.description,
			0 AS likes,
			0 AS comments,
			FALSE AS is_liked,
			tk.created_at
		FROM tasks tk
		JOIN users u
			ON u.user_id = tk.user_id

	) AS feed

	ORDER BY RANDOM()
	LIMIT $2
	`

	ctx, cancel := context.WithTimeout(
		context.Background(),
		3*time.Second,
	)
	defer cancel()

	rows, err := api.DB.QueryContext(
		ctx,
		query,
		userID,
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
			&feed.Likes,
			&feed.Comments,
			&feed.IsLiked,
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
