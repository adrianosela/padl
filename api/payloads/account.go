package payloads

import "errors"

// RegistrationRequest contains input for user registration
type RegistrationRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	PubKey   string `json:"public_key"`
}

// LoginRequest contains input for user login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse contains the response to a login request
type LoginResponse struct {
	Token string `json:"token"`
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

// Validate validates a login request payload
func (l *LoginRequest) Validate() error {
	if l.Email == "" {
		return errors.New("no email provided")
	}
	if l.Password == "" {
		return errors.New("no password provided")
	}
	return nil
}
