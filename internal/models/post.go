package models

import "time"

type SocialPost struct {
	PostID      int
	UserID      int
	Title       string
	Description string

	CreatedAt time.Time
}
type Like struct {
	UserID    int
	PostID    int
	CreatedAt time.Time
}

type Comment struct {
	CommentID  int       `json:"comment_id"`
	Post_id    int       `json:"post_id"`
	UserID     int       `json:"user_id"`
	Comment    string    `json:"comment"`
	CreatedAt  time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
}

type PostAuthor struct {
	PostID          int
	UserID          int
	Username        string
	DiscordUsername string
	Avatar          string
}

type PostResponse struct {
	PostID      int             `json:"post_id"`
	UserID      int             `json:"user_id"`
	Username    string          `json:"username"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Likes       int             `json:"likes"`
	Comments    int             `json:"comments"`
	IsLiked     bool            `json:"is_liked"`
	CommentList []CommentResponse `json:"comment_list"`
	CreatedAt   time.Time       `json:"created_at"`
}


type CommentResponse struct {
	CommentID      int       `json:"comment_id"`
	UserID         int       `json:"user_id"`
	Username       string    `json:"username"`
	Avatar         string    `json:"avatar"`
	Comment        string    `json:"comment"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}