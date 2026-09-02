package models

type Register struct {
    DiscordUserName string   `json:"DiscordUserName" validate:"required,min=3,max=50"`
    UserName        string   `json:"Username" validate:"required,min=3,max=50"`
    Roles           []string `json:"role" validate:"required,min=1"`
    Password        string   `json:"password" validate:"required,min=8"`
}

type Login struct {
	User_id  int    `json:"user_id"`
	Username string `validate:"required,min=3,max=50"`
	Password string `validate:"required,min=8"`
	Role     string `json:"role"`
}

type LoginTokens struct {
	AccessToken  string
	RefreshToken string
}
