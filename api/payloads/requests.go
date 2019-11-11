package payloads

// RegistrationRequest contains input for user registration
type RegistrationRequest struct {
	Email  string `json:"email"`
	PubKey string `json:"pub"`
}
