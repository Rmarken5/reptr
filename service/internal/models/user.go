package models

type (
	AuthRequest struct {
		GrantType    string `json:"grant_type,omitempty"`
		Username     string `json:"username,omitempty"`
		Password     string `json:"password,omitempty"`
		Audience     string `json:"audience,omitempty"`
		ClientID     string `json:"client_id,omitempty"`
		ClientSecret string `json:"client_secret,omitempty"`
	}
	TokenResp struct {
		AccessToken string `json:"accessToken"`
		TokenType   string `json:"tokenType"`
		ExpiresAt   int    `json:"expiresAt"`
	}
	User struct {
		ID       string `bson:"_id"`
		Username string `bson:"username"`
	}
)

type UserIDCtxKey struct {
}

var UserIDKey = UserIDCtxKey{}
