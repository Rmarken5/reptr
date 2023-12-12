package models

type (
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
