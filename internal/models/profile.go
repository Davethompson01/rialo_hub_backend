package models

type Profile struct {
    UserID           int             `json:"user_id"`
    Profile_pics     string          `json:"profile_pics"`
    Discord_username string          `json:"discord_username"`
    Username         string          `json:"username"`
    Role             string          `json:"role"`

    Tasks            []Task          `json:"tasks"`
    Posts            []PostResponse  `json:"posts"`
}