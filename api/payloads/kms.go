package payloads

import "errors"

// DecryptSecretRequest TODO
type DecryptSecretRequest struct {
	Secret string `json:"secret"`
}

// DecryptSecretResponse TODO
type DecryptSecretResponse struct {
	Message string `json:"message"`
}

// Validate TODO
func (a *DecryptSecretRequest) Validate() error {
	if a.Secret == "" {
		return errors.New("no secret provided")
	}
	return nil
}
