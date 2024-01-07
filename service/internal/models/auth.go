package models

type (
	RegistrationRequest struct {
		ClientID   string `json:"client_id,omitempty"`
		Email      string `json:"email,omitempty"`
		Password   string `json:"password,omitempty"`
		Connection string `json:"connection,omitempty"`
		Username   string `json:"username,omitempty"`
	}
	RegistrationUser struct {
		ID            string `json:"_id"`
		Email         string `json:"email"`
		EmailVerified bool   `json:"email_verified"`
	}

	RegistrationError struct {
		Name        string `json:"name"`
		Code        string `json:"code"`
		Description string `json:"description"`
		StatusCode  int    `json:"statusCode"`
	}
)

func (ru RegistrationUser) IsZero() bool {
	return ru.Email == "" && ru.ID == ""
}
func (re RegistrationError) IsZero() bool {
	return re.Name == "" && re.Code == "" && re.Description == "" && re.StatusCode == 0
}