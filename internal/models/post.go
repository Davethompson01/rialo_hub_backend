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

type PostsStatusMe struct {
	Employer_id  int  `json:"employer_id"`
	Created_Post bool `json:"created_post"`
	Status       string
}


type PostAuthor struct {
	ID          int    `json:"id"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	Avatar      string `json:"avatar"`
}

type PostResponse struct {
	ID          int        `json:"id"`
	Author      PostAuthor `json:"author"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Likes       int        `json:"likes"`
	Comments    int        `json:"comments"`
	IsLiked     bool       `json:"is_liked"`
	CreatedAt   time.Time  `json:"created_at"`
}