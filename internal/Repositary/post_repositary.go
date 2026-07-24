package repositary

import (
	"context"
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


// func LikePost(api *config.ApiConfig, post models.Like){
// query := ``
// }


