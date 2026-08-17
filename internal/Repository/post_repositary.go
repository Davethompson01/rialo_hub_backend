package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Davethompson01/rialo_hub_backend/config"
	"github.com/Davethompson01/rialo_hub_backend/internal/models"
)

func CreatePost(api *config.ApiConfig, post models.SocialPost) (int, error) {
	query := `
		INSERT INTO socialposts (
			user_id,
			title,
			description,
			created_at
		)
		VALUES ($1, $2, $3, $4)
		RETURNING post_id
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var postID int

	err := api.DB.QueryRowContext(
		ctx,
		query,
		post.UserID,
		post.Title,
		post.Description,
		post.CreatedAt,
	).Scan(&postID)

	if err != nil {
		return 0, err
	}

	return postID, nil
}

func BatchInsertLikes(api *config.ApiConfig, likes []models.Like) error {
	if len(likes) == 0 {
		return nil
	}

	var (
		args         []interface{}
		placeholders []string
	)

	for i, like := range likes {
		placeholders = append(
			placeholders,
			fmt.Sprintf("($%d, $%d, $%d)", i*3+1, i*3+2, i*3+3),
		)

		args = append(args,
			like.UserID,
			like.PostID,
			like.CreatedAt,
		)
	}

	query := `
		INSERT INTO likes (
			user_id,
			post_id,
			created_at
		)
		VALUES ` + strings.Join(placeholders, ",")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := api.DB.ExecContext(ctx, query, args...)

	return err
}

func CreateComment(
	api *config.ApiConfig,
	comment models.Comment,
) (int, error) {

	query := `
		INSERT INTO comments (
			post_id,
			user_id,
			comment,
			created_at,
			updated_at
		)
		VALUES (
			$1,
			$2,
			$3,
			CURRENT_TIMESTAMP,
			CURRENT_TIMESTAMP
		)
		RETURNING comment_id
	`

	var commentID int

	err := api.DB.QueryRow(
		query,
		comment.Post_id,
		comment.UserID,
		comment.Comment,
	).Scan(&commentID)

	if err != nil {
		return 0, err
	}

	return commentID, nil
}
func PostAuthor(api *config.ApiConfig, postID int) (models.PostAuthor, error) {

	var author models.PostAuthor

	query := `
		SELECT
			sp.post_id,
			u.user_id,
			u.username,
			u.discord_username,
			u.profile_pics
		FROM socialposts sp
		JOIN users u
			ON u.user_id = sp.user_id
		WHERE sp.post_id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := api.DB.QueryRowContext(ctx, query, postID).Scan(
		&author.PostID,
		&author.UserID,
		&author.Username,
		&author.DiscordUsername,
		&author.Avatar,
	)
	if err != nil {
		return models.PostAuthor{}, err
	}

	return author, nil
}

func PostResponse(api *config.ApiConfig, postID, userID int) (models.PostResponse, error) {
	var post models.PostResponse

	query := `
		SELECT
    sp.post_id,
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
          AND l.user_id = $2
    ) AS is_liked,
    sp.created_at
FROM socialposts sp
JOIN users u
    ON u.user_id = sp.user_id
WHERE sp.post_id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := api.DB.QueryRowContext(
		ctx,
		query,
		postID,
		userID,
	).Scan(
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
		return models.PostResponse{}, err
	}

	return post, nil
}

func UpdatePost(api *config.ApiConfig, post models.SocialPost) error {
	query := `
		UPDATE socialposts
		SET
			title = $1,
			description = $2
		WHERE post_id = $3
	`

	_, err := api.DB.Exec(
		query,
		post.Title,
		post.Description,
		post.PostID,
	)

	return err
}

func DeletePost(api *config.ApiConfig, postID int) error {
	query := `
		DELETE FROM socialposts
		WHERE post_id = $1
	`

	_, err := api.DB.Exec(query, postID)

	return err
}

func LikePost(api *config.ApiConfig, postID, userID int) error {
	tx, err := api.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback()

	var exists bool

	err = tx.QueryRow(`
		SELECT EXISTS (
			SELECT 1
			FROM socialposts
			WHERE post_id = $1
		)
	`, postID).Scan(&exists)

	if err != nil {
		return fmt.Errorf("failed to check post: %w", err)
	}

	if !exists {
		return errors.New("post not found")
	}

	result, err := tx.Exec(`
		INSERT INTO likes (
			user_id,
			post_id,
			created_at
		)
		VALUES ($1, $2, CURRENT_TIMESTAMP)
		ON CONFLICT (user_id, post_id) DO NOTHING
	`, userID, postID)

	if err != nil {
		return fmt.Errorf("failed to insert like: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected > 0 {
		_, err = tx.Exec(`
			UPDATE socialposts
			SET likes = likes + 1
			WHERE post_id = $1
		`, postID)

		if err != nil {
			return fmt.Errorf("failed to update like count: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit like: %w", err)
	}

	return nil
}

func UnlikePost(api *config.ApiConfig, postID, userID int) error {
	tx, err := api.DB.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	result, err := tx.Exec(`
		DELETE FROM likes
		WHERE user_id = $1
		  AND post_id = $2
	`, userID, postID)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	// Only decrement if a like actually existed.
	if rowsAffected > 0 {
		_, err = tx.Exec(`
			UPDATE socialposts
			SET likes = GREATEST(likes - 1, 0)
			WHERE post_id = $1
		`, postID)

		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func UpdateComment(
	api *config.ApiConfig,
	commentID int,
	content string,
) error {

	query := `
		UPDATE comments
		SET
			comment = $1,
			updated_at = CURRENT_TIMESTAMP
		WHERE comment_id = $2
	`

	_, err := api.DB.Exec(
		query,
		content,
		commentID,
	)

	return err
}

func IsPostOwner(api *config.ApiConfig, postID, userID int) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM socialposts
			WHERE post_id = $1
			  AND user_id = $2
		)
	`

	var exists bool

	err := api.DB.QueryRow(
		query,
		postID,
		userID,
	).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil
}

// func PostExists(api *config.ApiConfig, postID int) (bool, error) {
// 	query := `
// 		SELECT EXISTS (
// 			SELECT 1
// 			FROM socialposts
// 			WHERE post_id = $1
// 		)
// 	`

// 	var exists bool

// 	err := api.DB.QueryRow(
// 		query,
// 		postID,
// 	).Scan(&exists)

// 	if err != nil {
// 		return false, err
// 	}

// 	return exists, nil
// }



func PostExists(api *config.ApiConfig, postID int) (bool, error) {
	var exists bool

	err := api.DB.QueryRow(`
		SELECT EXISTS (
			SELECT 1
			FROM socialposts
			WHERE post_id = $1
		)
	`, postID).Scan(&exists)

	if err != nil {
		return false, err
	}

	fmt.Println("========== POST EXISTS ==========")
	fmt.Println("postID:", postID)
	fmt.Println("exists:", exists)
	fmt.Println("=================================")

	return exists, nil
}
func IsCommentOwner(api *config.ApiConfig, commentID, userID int) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM comments
			WHERE comment_id = $1
			  AND user_id = $2
		)
	`

	var exists bool

	err := api.DB.QueryRow(
		query,
		commentID,
		userID,
	).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil
}

func DeleteComment(api *config.ApiConfig, commentID int) error {
	query := `
		DELETE FROM comments
		WHERE comment_id = $1
	`

	_, err := api.DB.Exec(
		query,
		commentID,
	)

	return err
}
func GetPostFeed(
	api *config.ApiConfig,
	userID int,
	limit int,
	offset int,
) ([]models.PostResponse, error) {

	query := `
		SELECT
			sp.post_id,
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
		ORDER BY sp.created_at DESC
		LIMIT $2
		OFFSET $3
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
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.PostResponse

	for rows.Next() {
		var post models.PostResponse

		err := rows.Scan(
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
			return nil, err
		}

		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func GetPostComments(
	api *config.ApiConfig,
	postID int,
) ([]models.CommentResponse, error) {

	query := `
		SELECT
			c.comment_id,
			u.user_id,
			u.username,
			COALESCE(u.profile_pics, '') AS profile_pics,
			c.comment,
			c.created_at,
			c.updated_at
		FROM comments c
		JOIN users u
			ON u.user_id = c.user_id
		WHERE c.post_id = $1
		ORDER BY c.created_at ASC
	`

	ctx, cancel := context.WithTimeout(
		context.Background(),
		3*time.Second,
	)
	defer cancel()

	rows, err := api.DB.QueryContext(
		ctx,
		query,
		postID,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var comments []models.CommentResponse

	for rows.Next() {
		var comment models.CommentResponse

		err := rows.Scan(
			&comment.CommentID,
			&comment.UserID,
			&comment.Username,
			&comment.Avatar,
			&comment.Comment,
			&comment.CreatedAt,
			&comment.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}
