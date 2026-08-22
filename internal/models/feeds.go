package models

import "time"

type TrendingPost struct {
	Post PostResponse `json:"post"`
}

type FeedItem struct {
	Post      PostResponse `json:"post"`
	Following bool         `json:"following"`
}

type DashboardFeed struct {
	ID          int       `json:"id"`
	Type        string    `json:"type"`
	UserID      int       `json:"user_id"`
	Username    string    `json:"username"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Likes       int       `json:"likes"`
	Comments    int       `json:"comments"`
	IsLiked     bool      `json:"is_liked"`
	CreatedAt   time.Time `json:"created_at"`
}