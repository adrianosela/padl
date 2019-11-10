package payloads

type RegistrationRequest struct {
	Email  string `json:"email"`
	PubKey string `json:"pub"`
}
