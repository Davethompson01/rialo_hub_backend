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

func GetUserProfile(api *config.ApiConfig, userID int) (models.Profile, error) {
	profile := models.Profile{
		Tasks: []models.Task{},
		Posts: []models.PostResponse{},
	}

	// Get user
	userQuery := `
		SELECT
			user_id,
			COALESCE(profile_pics, ''),
			COALESCE(discord_username, ''),
			COALESCE(username, ''),
			role
		FROM users
		WHERE user_id = $1
	`

	err := api.DB.QueryRow(
		userQuery,
		userID,
	).Scan(
		&profile.UserID,
		&profile.Profile_pics,
		&profile.Discord_username,
		&profile.Username,
		&profile.Role,
	)

	if err != nil {
		return models.Profile{}, err
	}

	// Get user's tasks
	taskQuery := `
		SELECT
			task_id,
			user_id,
			title,
			description,
			reward,
		
			status,
			deadline,
			created_at
		FROM tasks
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := api.DB.Query(taskQuery, userID)
	if err != nil {
		return models.Profile{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var task models.Task

		err := rows.Scan(
			&task.ID,
			&task.UserID,
			&task.Title,
			&task.Description,
			&task.Reward,
			&task.Status,
			&task.Deadline,
			&task.CreatedAt,
		)

		if err != nil {
			return models.Profile{}, err
		}

		profile.Tasks = append(profile.Tasks, task)
	}

	if err := rows.Err(); err != nil {
		return models.Profile{}, err
	}

	// Get user's posts
	postQuery := `
		SELECT
			p.post_id,
			p.user_id,
			u.username,
			p.title,
			p.description,
			p.likes,
			p.comments,
			CASE
				WHEN l.user_id IS NOT NULL THEN true
				ELSE false
			END AS is_liked,
			p.created_at
		FROM socialposts p
		JOIN users u
			ON u.user_id = p.user_id
		LEFT JOIN likes l
			ON l.post_id = p.post_id
			AND l.user_id = $2
		WHERE p.user_id = $1
		ORDER BY p.created_at DESC
	`

	// $1 = profile user
	// $2 = current logged-in user
	postRows, err := api.DB.Query(
		postQuery,
		userID,
		userID,
	)

	if err != nil {
		return models.Profile{}, err
	}
	defer postRows.Close()

	for postRows.Next() {
		var post models.PostResponse

		err := postRows.Scan(
			&post.PostID,
			&post.UserID,
			&post.Username,
			&post.Title,
			&post.Description,
			&post.Likes,
			&post.Comments,
			&post.IsLiked,
			&post.CreatedAt,
		)

		if err != nil {
			return models.Profile{}, err
		}

		post.CommentList = []models.CommentResponse{}

		profile.Posts = append(profile.Posts, post)
	}

	if err := postRows.Err(); err != nil {
		return models.Profile{}, err
	}

	return profile, nil
}
