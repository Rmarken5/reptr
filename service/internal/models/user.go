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

// SubjectCtxKey to be used with auth provider values
type subjectCtxKey struct {
}
type userNameCtxKey struct {
}

// SubjectKey use to get and set auth provider values on context
var SubjectKey = subjectCtxKey{}

// UserNameKey use to get and set username on context
var UserNameKey = userNameCtxKey{}
