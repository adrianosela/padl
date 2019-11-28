package payloads

import "errors"

// DecryptSecretRequest contains secert to decrypt
type DecryptSecretRequest struct {
	Secret string `json:"secret"`
}

// DecryptSecretResponse contains reponse message to a secret decryption request
type DecryptSecretResponse struct {
	Message string `json:"message"`
}

// Validate validates a secret decryption request
func (a *DecryptSecretRequest) Validate() error {
	if a.Secret == "" {
		return errors.New("no secret provided")
	}
	return nil
}
