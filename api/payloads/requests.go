package payloads

import "errors"

// RegistrationRequest contains input for user registration
type RegistrationRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	PubKey   string `json:"public_key"`
}

// Validate validates a registration request payload
func (r *RegistrationRequest) Validate() error {
	if r.Email == "" {
		return errors.New("no email provided")
	}
	if r.Password == "" {
		return errors.New("no password provided")
	}
	// TODO: check PW complex enough
	if r.PubKey == "" {
		return errors.New("no public key provided")
	}
	// TODO: check pub key is valid RSA and at least 2048 bit
	return nil
}
