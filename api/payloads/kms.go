package payloads

import "errors"

// AddUserToKeyRequest TODO
type AddUserToKeyRequest struct {
	Email        string `json:"email"`
	PrivilegeLvl int    `json:"privilege"`
}

// DecryptSecretRequest TODO
type DecryptSecretRequest struct {
	Secret string `json:"secret"`
}

// RemoveUserFromKeyRequest TODO
type RemoveUserFromKeyRequest struct {
	Email string `json:"email"`
}

// DecryptSecretResponse TODO
type DecryptSecretResponse struct {
	Message string `json:"message"`
}

// Validate TODO
func (a *AddUserToKeyRequest) Validate() error {
	if a.Email == "" {
		return errors.New("no email provided")
	}
	if a.PrivilegeLvl < 0 || a.PrivilegeLvl > 3 {
		return errors.New("invalid privilege level provided")
	}
	return nil
}

// Validate TODO
func (a *RemoveUserFromKeyRequest) Validate() error {
	if a.Email == "" {
		return errors.New("no email provided")
	}
	return nil
}

// Validate TODO
func (a *DecryptSecretRequest) Validate() error {
	if a.Secret == "" {
		return errors.New("no secret provided")
	}
	return nil
}
