package models

type TrendingPost struct {
	Post PostResponse `json:"post"`
}

type FeedItem struct {
	Post      PostResponse `json:"post"`
	Following bool         `json:"following"`
}
