package models

type Register struct {
	DiscordUserName string `validate:"required,min=3,max=50"`
	UserName        string `validate:"required,min=3,max=50"`
	Role            string `json:"role"`
	Password        string `validate:"required,min=8"`
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
