package models

type AccessDetails struct {
	TokenUuid string
	UserId    string
	UserName  string
}

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	TokenUuid    string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

type User struct {
	Id       string `json:"id"`
	Name     string `json:"name" `
	Email    string `json:"email" `
	Role     string `json:"role"`
	Password string `json:"password"`
}
