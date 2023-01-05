package model

// TokenRes represents data about register user request
type TokenRes struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
type AccessTokenRes struct {
	AccessToken string `json:"access_token"`
}
