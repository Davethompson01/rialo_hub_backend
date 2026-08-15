package services

import (
	"errors"
	"fmt"

	"github.com/Davethompson01/rialo_hub_backend/config"
	repository "github.com/Davethompson01/rialo_hub_backend/internal/Repository"
	"github.com/Davethompson01/rialo_hub_backend/internal/models"
	"github.com/Davethompson01/rialo_hub_backend/internal/validation"
)

func CreatePost(api *config.ApiConfig, post models.SocialPost) (models.SocialPost, error) {
	if err := validation.ValidatePost(post); err != nil {
		return models.SocialPost{}, err
	}

	postID, err := repository.CreatePost(api, post)
	if err != nil {
		return models.SocialPost{}, err
	}

	post.PostID = postID

	return post, nil
}

func GetPost(
	api *config.ApiConfig,
	postID int,
	userID int,
) (models.PostResponse, error) {

	post, err := repository.PostResponse(api, postID, userID)
	if err != nil {
		return models.PostResponse{}, fmt.Errorf("failed to get post: %w", err)
	}

	return post, nil
}

func UpdatePost(
	api *config.ApiConfig,
	post models.SocialPost,
	userID int,
) error {

	ownsPost, err := repository.IsPostOwner(
		api,
		post.PostID,
		userID,
	)
	if err != nil {
		return fmt.Errorf("failed to verify post ownership: %w", err)
	}

	if !ownsPost {
		return errors.New("you are not allowed to update this post")
	}

	if err := validation.ValidatePost(post); err != nil {
		return err
	}

	if err := repository.UpdatePost(api, post); err != nil {
		return fmt.Errorf("failed to update post: %w", err)
	}

	return nil
}

func DeletePost(
	api *config.ApiConfig,
	postID int,
	userID int,
) error {

	ownsPost, err := repository.IsPostOwner(
		api,
		postID,
		userID,
	)
	if err != nil {
		return fmt.Errorf("failed to verify post ownership: %w", err)
	}

	if !ownsPost {
		return errors.New("you are not allowed to delete this post")
	}

	if err := repository.DeletePost(api, postID); err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}

	return nil
}

func LikePost(
	api *config.ApiConfig,
	postID int,
	userID int,
) error {

	exists, err := repository.PostExists(api, postID)
	if err != nil {
		return fmt.Errorf("failed to check post: %w", err)
	}

	if !exists {
		return errors.New("post not found")
	}

	if err := repository.LikePost(api, userID, postID); err != nil {
		return fmt.Errorf("failed to like post: %w", err)
	}

	return nil
}

func UnlikePost(
	api *config.ApiConfig,
	postID int,
	userID int,
) error {

	if err := repository.UnlikePost(api, userID, postID); err != nil {
		return fmt.Errorf("failed to unlike post: %w", err)
	}

	return nil
}

func CreateComment(
	api *config.ApiConfig,
	comment models.Comment,
) (models.Comment, error) {

	if err := validation.ValidateComment(comment); err != nil {
		return models.Comment{}, err
	}

	commentID, err := repository.CreateComment(api, comment)
	if err != nil {
		return models.Comment{}, fmt.Errorf(
			"failed to create comment: %w",
			err,
		)
	}

	comment.CommentID = commentID

	return comment, nil
}

func UpdateComment(
	api *config.ApiConfig,
	commentID int,
	userID int,
	content string,
) error {

	ownsComment, err := repository.IsCommentOwner(
		api,
		commentID,
		userID,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to verify comment ownership: %w",
			err,
		)
	}

	if !ownsComment {
		return errors.New(
			"you are not allowed to edit this comment",
		)
	}

	if err := repository.UpdateComment(
		api,
		commentID,
		content,
	); err != nil {
		return fmt.Errorf(
			"failed to update comment: %w",
			err,
		)
	}

	return nil
}

func DeleteComment(
	api *config.ApiConfig,
	commentID int,
	userID int,
) error {

	ownsComment, err := repository.IsCommentOwner(
		api,
		commentID,
		userID,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to verify comment ownership: %w",
			err,
		)
	}

	if !ownsComment {
		return errors.New(
			"you are not allowed to delete this comment",
		)
	}

	if err := repository.DeleteComment(
		api,
		commentID,
	); err != nil {
		return fmt.Errorf(
			"failed to delete comment: %w",
			err,
		)
	}

	return nil
}

func GetPostFeed(
	api *config.ApiConfig,
	userID int,
	limit int,
	offset int,
) ([]models.PostResponse, error) {

	posts, err := repository.GetPostFeed(
		api,
		userID,
		limit,
		offset,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get post feed: %w",
			err,
		)
	}

	return posts, nil
}

func GetPostComments(
	api *config.ApiConfig,
	postID int,
) ([]models.CommentResponse, error) {

	exists, err := repository.PostExists(api, postID)
	if err != nil {
		return nil, fmt.Errorf("failed to check post: %w", err)
	}

	if !exists {
		return nil, errors.New("post not found")
	}

	comments, err := repository.GetPostComments(api, postID)
	if err != nil {
		return nil, fmt.Errorf("failed to get post comments: %w", err)
	}

	return comments, nil
}
