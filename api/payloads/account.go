package payloads

import (
	"errors"

	"github.com/adrianosela/padl/lib/keys"
)

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

// RotateKeyRequest contains input for key rotation
type RotateKeyRequest struct {
	PubKey string `json:"public_key"`
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
	if r.PubKey == "" {
		return errors.New("no public key provided")
	}
	if _, err := keys.DecodePubKeyPEM([]byte(r.PubKey)); err != nil {
		return errors.New("provided key is not a PEM encoded RSA public key")
	}
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

// Validate validates a key rotation request
func (r *RotateKeyRequest) Validate() error {
	if r.PubKey == "" {
		return errors.New("no public key provided")
	}
	if _, err := keys.DecodePubKeyPEM([]byte(r.PubKey)); err != nil {
		return errors.New("provided key is not a PEM encoded RSA public key")
	}
	return nil
}
