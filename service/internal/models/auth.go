package models

type RegistrationRequest struct {
	ClientID   string `json:"client_id,omitempty"`
	Email      string `json:"email,omitempty"`
	Password   string `json:"password,omitempty"`
	Connection string `json:"connection,omitempty"`
	Username   string `json:"username,omitempty"`
}
