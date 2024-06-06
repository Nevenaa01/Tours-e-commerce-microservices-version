package model

type AuthenticationTokensDto struct {
	ID          int    `json:"id"`
	AccessToken string `json:"access_token"`
}
