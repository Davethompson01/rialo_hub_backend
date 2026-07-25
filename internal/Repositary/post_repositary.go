package repositary

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Davethompson01/rialo_hub_backend/config"
	"github.com/Davethompson01/rialo_hub_backend/internal/models"
)

func CreatePost(api *config.ApiConfig, post models.SocialPost) error {
	query := `INSERT INTO socialposts(user_id, title, description, created_at) 
	VALUES(($1, $2, $3, $4)`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := api.DB.ExecContext(
		ctx,
		query,
		post.UserID,
		post.Title,
		post.Title,
		post.CreatedAt,
	)
	return err

}

func BatchInsertLikes(api *config.ApiConfig, likes []models.Like) error {
	if len(likes) == 0 {
		return nil
	}

	var (
		args         []interface{}
		placeholders []string
	)

	query := `
		INSERT INTO likes(user_id, post_id, created_at)
		VALUES
	`

	for i, like := range likes {
		// Each row has 3 columns
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

	query += strings.Join(placeholders, ",")

	_, err := api.DB.Exec(query, args...)
	return err
}

func SocialFeedsComment(api *config.ApiConfig, comment models.Comment) error {
	query := `
		INSERT INTO comments (
			post_id,
			user_id,
			comment,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5)
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := api.DB.ExecContext(
		ctx,
		query,
		comment.Post_id,
		comment.UserID,
		comment.Comment,
		comment.CreatedAt,
		comment.Updated_at,
	)

	return err
}
func PostAuthor(api *config.ApiConfig, postID int) (models.PostAuthor, error) {

	var author models.PostAuthor

	query := `
		SELECT
			sp.post_id,
			u.user_id,
			u.username,
			u.discord_username,
			u.avatar
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

func PostResponse(api *config.ApiConfig, postID int) (models.PostResponse, error) {

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
			sp.created_at
		FROM socialposts sp
		JOIN users u
			ON u.user_id = sp.user_id
		WHERE sp.post_id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := api.DB.QueryRowContext(ctx, query, postID).Scan(
		&post.PostID,
		&post.UserID,
		&post.Username,
		&post.Title,
		&post.Description,
		&post.Likes,
		&post.Comments,
		&post.CreatedAt,
	)
	if err != nil {
		return models.PostResponse{}, err
	}

	return post, nil
}