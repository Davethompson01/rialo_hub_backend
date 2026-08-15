package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	middleware "github.com/Davethompson01/rialo_hub_backend/Middleware"
	"github.com/Davethompson01/rialo_hub_backend/config"
	auth "github.com/Davethompson01/rialo_hub_backend/internal/Auth"
	"github.com/Davethompson01/rialo_hub_backend/internal/models"
	"github.com/Davethompson01/rialo_hub_backend/internal/services"
	"github.com/go-chi/chi"
)

func CreatePost(api *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var post models.SocialPost
		if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
			RespondWithJson(
				w,
				http.StatusBadRequest,
				false,
				"Invalid request body",
				nil,
			)
			return
		}
		claims, ok := r.Context().Value(middleware.ClaimsKey).(*auth.Claims)
		if !ok || claims == nil {
			RespondWithJson(
				w,
				http.StatusUnauthorized,
				false,
				"Unauthorized",
				nil,
			)
			return
		}

		post.UserID = claims.UserID
		post, err := services.CreatePost(api, post)
		if err != nil {
			RespondWithJson(
				w,
				http.StatusBadRequest,
				false,
				err.Error(),
				nil,
			)
			return
		}

		RespondWithJson(
			w,
			http.StatusCreated,
			true,
			"Post Created Successfully",
			post,
		)
	}
}

func GetPost(api *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(middleware.ClaimsKey).(*auth.Claims)
		if !ok || claims == nil {
			RespondWithJson(
				w,
				http.StatusUnauthorized,
				false,
				"Unauthorized",
				nil,
			)
			return
		}

		postID, err := strconv.Atoi(chi.URLParam(r, "postID"))
		if err != nil {
			RespondWithJson(
				w,
				http.StatusBadRequest,
				false,
				"Invalid post ID",
				nil,
			)
			return
		}

		post, err := services.GetPost(
			api,
			postID,
			claims.UserID,
		)
		if err != nil {
			RespondWithJson(
				w,
				http.StatusNotFound,
				false,
				err.Error(),
				nil,
			)
			return
		}

		RespondWithJson(
			w,
			http.StatusOK,
			true,
			"Post retrieved successfully",
			post,
		)
	}
}

func UpdatePost(api *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(middleware.ClaimsKey).(*auth.Claims)
		if !ok || claims == nil {
			RespondWithJson(
				w,
				http.StatusUnauthorized,
				false,
				"Unauthorized",
				nil,
			)
			return
		}

		postID, err := strconv.Atoi(chi.URLParam(r, "postID"))
		if err != nil {
			RespondWithJson(
				w,
				http.StatusBadRequest,
				false,
				"Invalid post ID",
				nil,
			)
			return
		}

		var post models.SocialPost

		if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
			RespondWithJson(
				w,
				http.StatusBadRequest,
				false,
				"Invalid request body",
				nil,
			)
			return
		}

		post.PostID = postID

		err = services.UpdatePost(
			api,
			post,
			claims.UserID,
		)
		if err != nil {
			RespondWithJson(
				w,
				http.StatusBadRequest,
				false,
				err.Error(),
				nil,
			)
			return
		}

		RespondWithJson(
			w,
			http.StatusOK,
			true,
			"Post updated successfully",
			nil,
		)
	}
}

func DeletePost(api *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(middleware.ClaimsKey).(*auth.Claims)
		if !ok || claims == nil {
			RespondWithJson(
				w,
				http.StatusUnauthorized,
				false,
				"Unauthorized",
				nil,
			)
			return
		}

		postID, err := strconv.Atoi(chi.URLParam(r, "postID"))
		if err != nil {
			RespondWithJson(
				w,
				http.StatusBadRequest,
				false,
				"Invalid post ID",
				nil,
			)
			return
		}

		err = services.DeletePost(
			api,
			postID,
			claims.UserID,
		)
		if err != nil {
			RespondWithJson(
				w,
				http.StatusBadRequest,
				false,
				err.Error(),
				nil,
			)
			return
		}

		RespondWithJson(
			w,
			http.StatusOK,
			true,
			"Post deleted successfully",
			nil,
		)
	}
}

func LikePost(api *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(middleware.ClaimsKey).(*auth.Claims)
		if !ok || claims == nil {
			RespondWithJson(
				w,
				http.StatusUnauthorized,
				false,
				"Unauthorized",
				nil,
			)
			return
		}

		postID, err := strconv.Atoi(chi.URLParam(r, "postID"))
		if err != nil {
			RespondWithJson(
				w,
				http.StatusBadRequest,
				false,
				"Invalid post ID",
				nil,
			)
			return
		}

		err = services.LikePost(
			api,
			postID,
			claims.UserID,
		)
		if err != nil {
			RespondWithJson(
				w,
				http.StatusBadRequest,
				false,
				err.Error(),
				nil,
			)
			return
		}

		RespondWithJson(
			w,
			http.StatusOK,
			true,
			"Post liked successfully",
			nil,
		)
	}
}

func UnlikePost(api *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(middleware.ClaimsKey).(*auth.Claims)
		if !ok || claims == nil {
			RespondWithJson(
				w,
				http.StatusUnauthorized,
				false,
				"Unauthorized",
				nil,
			)
			return
		}

		postID, err := strconv.Atoi(chi.URLParam(r, "postID"))
		if err != nil {
			RespondWithJson(
				w,
				http.StatusBadRequest,
				false,
				"Invalid post ID",
				nil,
			)
			return
		}

		err = services.UnlikePost(
			api,
			postID,
			claims.UserID,
		)
		if err != nil {
			RespondWithJson(
				w,
				http.StatusBadRequest,
				false,
				err.Error(),
				nil,
			)
			return
		}

		RespondWithJson(
			w,
			http.StatusOK,
			true,
			"Post unliked successfully",
			nil,
		)
	}
}

func CreateComment(api *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(middleware.ClaimsKey).(*auth.Claims)
		if !ok || claims == nil {
			RespondWithJson(
				w,
				http.StatusUnauthorized,
				false,
				"Unauthorized",
				nil,
			)
			return
		}

		postID, err := strconv.Atoi(chi.URLParam(r, "postID"))
		if err != nil {
			RespondWithJson(
				w,
				http.StatusBadRequest,
				false,
				"Invalid post ID",
				nil,
			)
			return
		}

		var comment models.Comment

		if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
			RespondWithJson(
				w,
				http.StatusBadRequest,
				false,
				"Invalid request body",
				nil,
			)
			return
		}

		comment.Post_id = postID
		comment.UserID = claims.UserID

		createdComment, err := services.CreateComment(
			api,
			comment,
		)
		if err != nil {
			RespondWithJson(
				w,
				http.StatusBadRequest,
				false,
				err.Error(),
				nil,
			)
			return
		}

		RespondWithJson(
			w,
			http.StatusCreated,
			true,
			"Comment created successfully",
			createdComment,
		)
	}
}

func UpdateComment(api *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(middleware.ClaimsKey).(*auth.Claims)
		if !ok || claims == nil {
			RespondWithJson(
				w,
				http.StatusUnauthorized,
				false,
				"Unauthorized",
				nil,
			)
			return
		}

		commentID, err := strconv.Atoi(chi.URLParam(r, "commentID"))
		if err != nil {
			RespondWithJson(
				w,
				http.StatusBadRequest,
				false,
				"Invalid comment ID",
				nil,
			)
			return
		}

		var body struct {
			Comment string `json:"comment"`
		}

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			RespondWithJson(
				w,
				http.StatusBadRequest,
				false,
				"Invalid request body",
				nil,
			)
			return
		}

		err = services.UpdateComment(
			api,
			commentID,
			claims.UserID,
			body.Comment,
		)
		if err != nil {
			RespondWithJson(
				w,
				http.StatusBadRequest,
				false,
				err.Error(),
				nil,
			)
			return
		}

		RespondWithJson(
			w,
			http.StatusOK,
			true,
			"Comment updated successfully",
			nil,
		)
	}
}

func DeleteComment(api *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(middleware.ClaimsKey).(*auth.Claims)
		if !ok || claims == nil {
			RespondWithJson(
				w,
				http.StatusUnauthorized,
				false,
				"Unauthorized",
				nil,
			)
			return
		}

		commentID, err := strconv.Atoi(chi.URLParam(r, "commentID"))
		if err != nil {
			RespondWithJson(
				w,
				http.StatusBadRequest,
				false,
				"Invalid comment ID",
				nil,
			)
			return
		}

		err = services.DeleteComment(
			api,
			commentID,
			claims.UserID,
		)
		if err != nil {
			RespondWithJson(
				w,
				http.StatusBadRequest,
				false,
				err.Error(),
				nil,
			)
			return
		}

		RespondWithJson(
			w,
			http.StatusOK,
			true,
			"Comment deleted successfully",
			nil,
		)
	}
}


func GetPostFeed(api *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(middleware.ClaimsKey).(*auth.Claims)
		if !ok || claims == nil {
			RespondWithJson(
				w,
				http.StatusUnauthorized,
				false,
				"Unauthorized",
				nil,
			)
			return
		}

		page := 1
		limit := 20

		if value := r.URL.Query().Get("page"); value != "" {
			parsed, err := strconv.Atoi(value)
			if err != nil || parsed < 1 {
				RespondWithJson(
					w,
					http.StatusBadRequest,
					false,
					"Invalid page",
					nil,
				)
				return
			}

			page = parsed
		}

		if value := r.URL.Query().Get("limit"); value != "" {
			parsed, err := strconv.Atoi(value)
			if err != nil || parsed < 1 || parsed > 100 {
				RespondWithJson(
					w,
					http.StatusBadRequest,
					false,
					"Invalid limit",
					nil,
				)
				return
			}

			limit = parsed
		}

		offset := (page - 1) * limit

		posts, err := services.GetPostFeed(
			api,
			claims.UserID,
			limit,
			offset,
		)
		if err != nil {
			RespondWithJson(
				w,
				http.StatusInternalServerError,
				false,
				err.Error(),
				nil,
			)
			return
		}

		RespondWithJson(
			w,
			http.StatusOK,
			true,
			"Posts retrieved successfully",
			posts,
		)
	}
}



func GetPostComments(api *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Make sure the user is authenticated.
		claims, ok := r.Context().Value(middleware.ClaimsKey).(*auth.Claims)
		if !ok || claims == nil {
			RespondWithJson(
				w,
				http.StatusUnauthorized,
				false,
				"Unauthorized",
				nil,
			)
			return
		}

		
		postID, err := strconv.Atoi(
			chi.URLParam(r, "postID"),
		)

		if err != nil || postID <= 0 {
			RespondWithJson(
				w,
				http.StatusBadRequest,
				false,
				"Invalid post ID",
				nil,
			)
			return
		}

		// Get all comments belonging to the post.
		comments, err := services.GetPostComments(
			api,
			postID,
		)

		if err != nil {
			RespondWithJson(
				w,
				http.StatusNotFound,
				false,
				err.Error(),
				nil,
			)
			return
		}

		RespondWithJson(
			w,
			http.StatusOK,
			true,
			"Post comments retrieved successfully",
			comments,
		)
	}
}