package models

import "time"

type SocialPost struct {

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
	Post_id int `json:"post_id"`
	UserID  int `json:"user_id"`
	Comment    string `json:"comment"`
	CreatedAt  time.Time
	Updated_at time.Time
}

type PostAuthor struct {
	PostID           int
	UserID           int
	Username         string
	DiscordUsername  string
	Avatar           string
}

type PostResponse struct {
	PostID      int
	UserID      int
	Username    string
	Title       string
	Description string
	Likes       int
	Comments    int
	IsLiked     bool
	CreatedAt   time.Time
}