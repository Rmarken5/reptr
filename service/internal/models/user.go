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
		Username       string   `bson:"_id"`
		MemberOfGroups []string `bson:"member_of_groups"`
	}
)
